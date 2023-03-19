package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
)

const (
	recievedFile          = path + "/gildedrose_test.TestApprovalUpdateQuality.received.txt"
	approvedFile          = path + "/gildedrose_test.TestApprovalUpdateQuality.approved.txt"
	testCasesTemplateFile = "test_cases.yaml.tpl"
	testCasesOutputFile   = path + "/test_cases.yaml"
)

func cleanFile(fo *os.File) {
	if err := fo.Truncate(0); err != nil {
		fmt.Println(err)
	}
	fo.Seek(0, 0)
}

func updateApproved() {
	filename, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
	}

	var command []string

	command = []string{
		"rm",
		"-f",
		recievedFile,
		approvedFile,
	}

	runCommand(command, false)

	command = []string{
		"go",
		"test",
		filename,
	}

	runCommand(command, true)

	command = []string{
		"mv",
		recievedFile,
		approvedFile,
	}

	runCommand(command, false)
}

func verifyDuplicate(combos []Combo, c Combo) bool {
	for _, combo := range combos {
		if reflect.DeepEqual(combo, c) {
			return true
		}
	}
	return false
}

func GenTestCasesFile(allCombos []Combo, finalCombos []Combo, getValue func() float64) []Combo {
	t, err := template.ParseFiles(testCasesTemplateFile)
	if err != nil {
		panic(err)
	}

	// open output file
	fo, err := os.Create(testCasesOutputFile)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	var oldValue, value float64 = 0, 0

	comobosInitiliazed := len(finalCombos) != 0
	if comobosInitiliazed {
		if err = t.Execute(fo, finalCombos); err != nil {
			panic(err)
		}
		updateApproved()
		oldValue = getValue()
	}

	for _, c := range allCombos {
		if comobosInitiliazed && verifyDuplicate(finalCombos, c) {
			continue
		}
		tempCombos := append(finalCombos, c)
		if err = t.Execute(fo, tempCombos); err != nil {
			panic(err)
		}
		updateApproved()

		value = getValue()
		if oldValue < value {
			finalCombos = tempCombos
			oldValue = value
		}
		// time.Sleep(2 * time.Second)
		log.Println("----------------------> percent = ", value)
		cleanFile(fo)
		if value == 100 {
			break
		}
	}

	updateApproved()
	if err = t.Execute(fo, finalCombos); err != nil {
		panic(err)
	}

	return finalCombos
}
