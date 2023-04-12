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
	"sync"
	"text/template"
)

func updateApproved(p string) {

	utilities.DeleteFile(p + "/" + config.RecievedFile)
	utilities.DeleteFile(p + "/" + config.ApprovedFile)

	c := exec.Command("go", "test", p+"/"+config.Path)
	c.Dir = p
	utilities.RunCommand(c, true)

	utilities.MoveFile(p+"/"+config.Path+"/"+config.RecievedFile, p+"/"+config.Path+"/"+config.ApprovedFile)
}

func doSubCombosThing(combos *combo.Combos, customTest scores.CustomTest, t *template.Template, indexJobs <-chan int, arrayChannel *[]float64, wg *sync.WaitGroup, maxIndex *utilities.Safeindex) {
	defer wg.Done()

	for index := range indexJobs {
		subCombos := (*combos)[:index+1]

		if (*arrayChannel)[maxIndex.Value()] == 100 && index > maxIndex.Value() {
			log.Println("skip index ", index)
			return
		}

		h := subCombos.GetSubCombosHash()
		dir := utilities.CreateSandboxedEnv(h)

		// open output file
		fo, err := os.Create(dir + "/" + config.Path + "/" + config.TestCasesOutputFile)
		if err != nil {
			panic(err)
		}

		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()

		if err := t.Execute(fo, subCombos); err != nil {
			panic(err)
		}
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

	comobosInitiliazed := len(finalCombos) != 0
	if comobosInitiliazed {
		fmt.Println("comobosInitiliazed = true")
		// if err = t.Execute(fo, finalCombos); err != nil {
		// 	panic(err)
		// }
		// updateApproved(path)
		// oldValue = GetValue()
	}

	var wg sync.WaitGroup
	utilities.CleanDir(config.BaseTmpDir)
	arrayChannel := &[]float64{}

	// TODO: use channel
	maxIndex := utilities.Safeindex{}

	const numWorkers = 10
	// Create a buffered channel with a capacity of numWorkers
	indexJobs := make(chan int, numWorkers)

	// allCombosWithoutPointer := []combo.Combo{}
	// for _, c := range allCombos {
	// 	allCombosWithoutPointer = append(allCombosWithoutPointer, *c)
	// }

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go doSubCombosThing(&allCombos, customTest, t, indexJobs, arrayChannel, &wg, &maxIndex)
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

	// updateApproved(path)

	// open output file
	fo, err := os.Create(config.BasePath + "/" + config.Path + "/" + config.TestCasesOutputFile)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if err := t.Execute(fo, finalCombos); err != nil {
		panic(err)
	}

	return finalCombos
}
