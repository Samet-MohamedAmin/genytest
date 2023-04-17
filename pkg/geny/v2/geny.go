package geny

import (
	"fmt"
	"genytest/config"
	"genytest/pkg/combo"
	"genytest/scores"
	"genytest/utilities"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"sync"
	"text/template"
)

const (
	MAX_CONCURRENT = 8
	MAX_WORKERS    = 8
)

var (
	t          *template.Template
	customTest scores.CustomTest
)

type Score struct {
	value float64
	mutex sync.Mutex
}

func updateApproved(p string) {

	utilities.DeleteFile(path.Join(p, config.RecievedFile))
	utilities.DeleteFile(path.Join(p, config.ApprovedFile))

	c := exec.Command("go", "test", path.Join(p, config.Path))
	c.Dir = p
	utilities.RunCommand(c, true)

	err := os.Rename(path.Join(p, config.Path, config.RecievedFile), path.Join(p, config.Path, config.ApprovedFile))
	if err != nil {
		log.Println(err)
	}
}

func writeTemplate(t *template.Template, pathFile string, data any) {

	// open output file
	fo, err := os.Create(pathFile)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if err := t.Execute(fo, data); err != nil {
		panic(err)
	}
}

func setScore(
	combos combo.Combos,
	scores *[]Score,
	index int,
	jobNumbers chan int) {

	number := <-jobNumbers
	defer func() { jobNumbers <- number }()
	(*scores)[index].mutex.Lock()
	defer (*scores)[index].mutex.Unlock()
	if (*scores)[index].value != 0 {
		// log.Printf("Score[%d] already calculated\n", index)
		return
	}

	dir := utilities.CreateSandboxedEnv(fmt.Sprint(number))
	subCombos := combo.Combos{}
	for _, c := range combos[:index+1] {
		if !c.Useless {
			subCombos = append(subCombos, c)
		}
	}

	// open output file
	pathFile := path.Join(dir, config.Path, config.TestCasesOutputFile)
	writeTemplate(t, pathFile, subCombos)

	updateApproved(dir)

	(*scores)[index].value = customTest.GetValue(dir)
	log.Printf("scores[%d] = %.2f\n", index, (*scores)[index].value)

}

func printStat(scores *[]Score) {
	calculated := 0
	for i := 0; i < len(*scores); i++ {
		if (*scores)[i].value != 0 {
			calculated++
		}
	}
	percent := float64(calculated) * 100 / float64(len(*scores))
	log.Printf("--------------------------------------------> calculated %.2f%%\n", percent)
}

func markUseless(firstIndex, lastIndex int, combos combo.Combos, scores *[]Score) {
	if firstIndex > len(combos) {
		return
	}
	log.Printf("---------> marked usesless form %d to %d\n", firstIndex, lastIndex)
	for i := firstIndex; i <= lastIndex; i++ {
		(*scores)[i].value = (*scores)[lastIndex].value
		combos[i].Useless = true
	}
	printStat(scores)
}

func doTheSubThing(
	wg *sync.WaitGroup,
	combos combo.Combos,
	scores *[]Score,
	firstIndex int,
	lastIndex int,
	jobNumbers chan int,
	concurrent *utilities.Safeindex) {
	defer wg.Done()

	waitDepth := true
	if concurrent.ValueUnsafe() >= MAX_CONCURRENT {
		concurrent.Increase()
		waitDepth = concurrent.Value() >= MAX_CONCURRENT
	}

	medium := (firstIndex + lastIndex) / 2

	if !waitDepth {
		if lastIndex-firstIndex > 1 {
			wg.Add(2)
			go doTheSubThing(wg, combos, scores, firstIndex, medium, jobNumbers, concurrent)
			go doTheSubThing(wg, combos, scores, medium, lastIndex, jobNumbers, concurrent)
		}
		if lastIndex-firstIndex == 1 {
			setScore(combos, scores, firstIndex+1, jobNumbers)
			setScore(combos, scores, firstIndex, jobNumbers)
			setScore(combos, scores, lastIndex, jobNumbers)
			return
		}

	}

	if lastIndex-firstIndex <= 0 {
		return
	}

	setScore(combos, scores, firstIndex, jobNumbers)
	setScore(combos, scores, lastIndex, jobNumbers)
	if firstIndex == lastIndex {
		return
	}

	if (*scores)[firstIndex].value == 100 {
		markUseless(firstIndex+1, len(combos)-1, combos, scores)
	}

	if (*scores)[lastIndex].value == 100 && lastIndex != len(combos)-1 {
		markUseless(lastIndex+1, len(combos)-1, combos, scores)
	}

	if (*scores)[firstIndex].value == (*scores)[lastIndex].value {
		if (*scores)[firstIndex].value != 100 {
			markUseless(firstIndex+1, lastIndex, combos, scores)
		}
		return
	}

	if waitDepth {
		if lastIndex-firstIndex > 1 {
			wg.Add(2)
			go doTheSubThing(wg, combos, scores, firstIndex, medium, jobNumbers, concurrent)
			go doTheSubThing(wg, combos, scores, medium, lastIndex, jobNumbers, concurrent)
		}
		if lastIndex-firstIndex == 1 {
			setScore(combos, scores, firstIndex+1, jobNumbers)
			return
		}
	}
}

func setSingleScore(
	combos combo.Combos,
	scores *[]Score,
	wg *sync.WaitGroup,
	indexJobs <-chan int,
	jobNumber int) {
	defer wg.Done()
	dir := utilities.CreateSandboxedEnv(fmt.Sprint(jobNumber))
	for index := range indexJobs {
		// open output file
		pathFile := path.Join(dir, config.Path, config.TestCasesOutputFile)
		writeTemplate(t, pathFile, combos[index:index+1])

		updateApproved(dir)

		(*scores)[index].value = customTest.GetValue(dir)
		log.Printf("scores[%d] = %.2f\n", index, (*scores)[index].value)
	}

}

func setScoresFromStart(
	combos combo.Combos,
	scores *[]Score,
	wg *sync.WaitGroup,
	indexJobs <-chan int,
	jobNumber int) {
	defer wg.Done()
	dir := utilities.CreateSandboxedEnv(fmt.Sprint(jobNumber))
	for index := range indexJobs {
		// open output file
		pathFile := path.Join(dir, config.Path, config.TestCasesOutputFile)
		writeTemplate(t, pathFile, combos[0:index+1])

		updateApproved(dir)

		(*scores)[index].value = customTest.GetValue(dir)
		log.Printf("scores[%d] = %.2f\n", index, (*scores)[index].value)
	}

}

func sequentialFilter(combos combo.Combos) (finalCombos combo.Combos) {
	log.Println("------------------------------> start sequentialFilter")
	scores := []Score{}
	for i := 0; i < len(combos); i++ {
		scores = append(scores, Score{})
	}

	var wg sync.WaitGroup
	indexJobs := make(chan int, MAX_WORKERS)

	for i := 1; i <= MAX_WORKERS; i++ {
		wg.Add(1)
		go setSingleScore(combos, &scores, &wg, indexJobs, i)
	}

	for i := range combos {
		indexJobs <- i
	}

	close(indexJobs)
	wg.Wait()

	sort.Slice(combos, func(i, j int) bool {
		return scores[i].value > scores[j].value
	})

	scores = []Score{}
	for i := 0; i < len(combos); i++ {
		scores = append(scores, Score{})
	}
	indexJobs = make(chan int, MAX_WORKERS)

	for i := 1; i <= MAX_WORKERS; i++ {
		wg.Add(1)
		go setScoresFromStart(combos, &scores, &wg, indexJobs, i)
	}

	for i := range combos {
		indexJobs <- i
	}

	close(indexJobs)
	wg.Wait()

	oldResult := float64(0)
	for i := range combos {
		newResult := scores[i].value
		if oldResult < newResult {
			finalCombos = append(finalCombos, combos[i])
		}
		oldResult = newResult
	}
	log.Println("------------------------------> end sequentialFilter")

	return
}

func GenTestCasesFile(allCombos combo.Combos, finalCombos combo.Combos, ct scores.CustomTest) combo.Combos {
	customTest = ct
	if len(allCombos) < 2 {
		return allCombos
	}
	// read template
	var err error
	t, err = template.ParseFiles(config.TestCasesTemplateFile)
	if err != nil {
		panic(err)
	}

	// Create a buffered channel with a capacity of MAX_WORKERS
	var wg sync.WaitGroup
	jobNumbers := make(chan int, MAX_WORKERS)

	for i := 1; i <= MAX_WORKERS; i++ {
		jobNumbers <- i
	}

	scores := []Score{}
	for i := 0; i < len(allCombos); i++ {
		scores = append(scores, Score{})
	}
	wg.Add(1)

	concurrent := utilities.Safeindex{}
	concurrent.SetValue(0)

	go doTheSubThing(&wg, allCombos, &scores, 0, len(allCombos)-1, jobNumbers, &concurrent)

	wg.Wait()
	close(jobNumbers)

	mediumResults := combo.Combos{}
	oldResult := float64(0)
	for i := range allCombos {
		newResult := scores[i].value
		if oldResult < newResult {
			mediumResults = append(mediumResults, allCombos[i])
		}
		oldResult = newResult
	}

	finalCombos = append(finalCombos, sequentialFilter(mediumResults)...)

	outputPath := path.Join(config.BasePath, config.Path, config.TestCasesOutputFile)
	writeTemplate(t, outputPath, finalCombos)

	return finalCombos
}
