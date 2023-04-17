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
	"sync"
	"text/template"
)

func updateApproved(p string) {

	utilities.DeleteFile(path.Join(p, config.RecievedFile))
	utilities.DeleteFile(path.Join(p, config.ApprovedFile))

	c := exec.Command("go", "test", path.Join(p, config.Path))
	c.Dir = p
	utilities.RunCommand(c, true)

	os.Rename(path.Join(p, config.Path, config.RecievedFile), path.Join(p, config.Path, config.ApprovedFile))
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

func doSubCombosThing(
	combos *combo.Combos,
	customTest scores.CustomTest,
	t *template.Template,
	indexJobs <-chan int,
	arrayChannel *[]float64,
	wg *sync.WaitGroup,
	maxIndex *utilities.Safeindex,
	jobNumber int) {
	defer wg.Done()

	dir := utilities.CreateSandboxedEnv(fmt.Sprint(jobNumber))

	for index := range indexJobs {
		subCombos := (*combos)[:index+1]

		if (*arrayChannel)[maxIndex.Value()] == 100 && index > maxIndex.Value() {
			log.Println("skip index ", index)
			return
		}

		// open output file
		pathFile := path.Join(dir, config.Path, config.TestCasesOutputFile)
		writeTemplate(t, pathFile, subCombos)

		updateApproved(dir)

		if (*arrayChannel)[maxIndex.Value()] == 100 && index > maxIndex.Value() {
			log.Println("skip index ", index)
			return
		}
		value := customTest.GetValue(dir)
		if value > (*arrayChannel)[maxIndex.Value()] {
			maxIndex.SetValue(index)
			log.Printf("update maxValue = %f", value)
		}

		(*arrayChannel)[index] = value
	}
}

func GenTestCasesFile(allCombos combo.Combos, finalCombos combo.Combos, customTest scores.CustomTest) combo.Combos {
	if len(allCombos) < 2 {
		return allCombos
	}

	t, err := template.ParseFiles(config.TestCasesTemplateFile)
	if err != nil {
		panic(err)
	}

	// TODO: check if finalComobos is initiliazed

	var wg sync.WaitGroup
	utilities.CleanDir(config.BaseTmpDir)
	arrayChannel := &[]float64{}

	maxIndex := utilities.Safeindex{}

	const numWorkers = 6
	// Create a buffered channel with a capacity of numWorkers
	indexJobs := make(chan int, numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go doSubCombosThing(&allCombos, customTest, t, indexJobs, arrayChannel, &wg, &maxIndex, i)
	}

	// Send the jobs to the workers
	for i := range allCombos {
		*arrayChannel = append(*arrayChannel, 0)
		indexJobs <- i
		if (*arrayChannel)[maxIndex.Value()] == 100 {
			break
		}
	}

	// Close the jobs channel to signal that there are no more jobs
	close(indexJobs)

	wg.Wait()

	oldResult := float64(0)
	for i := range *arrayChannel {
		newResult := (*arrayChannel)[i]
		if oldResult < newResult {
			finalCombos = append(finalCombos, allCombos[i])
		}
		oldResult = newResult
	}

	outputPath := path.Join(config.BasePath, config.Path, config.TestCasesOutputFile)
	writeTemplate(t, outputPath, finalCombos)

	return finalCombos
}
