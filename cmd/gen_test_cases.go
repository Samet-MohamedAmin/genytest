package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	recievedFile = path + "/example_test.TestBla.received.txt"
	approvedFile = path + "/example_test.TestBla.approved.txt"
)

type combo map[string]string

type testTemplate struct {
	Id     string `yaml:"id"`
	Input  int    `yaml:"input"`
	Output int    `Yaml:"output"`
}

type possibleValue struct {
	Type           string   `yaml:"type"`
	PossibleValues []string `yaml:"possibleValues"`
}

type valuesByKey map[string]possibleValue

func GetPossibleValues() valuesByKey {
	filename, _ := filepath.Abs(possibleValuesFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var valuesMap valuesByKey

	err = yaml.Unmarshal(yamlFile, &valuesMap)
	if err != nil {
		panic(err)
	}

	return valuesMap
}

func getCombos(vmap valuesByKey) (output []combo) {
	for key, value := range vmap {
		for _, pv := range value.PossibleValues {
			item := map[string]string{
				key: pv,
			}
			output = append(output, item)
		}
	}
	return
}

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

func GenTestCasesFile(pv valuesByKey, getValue func() float64) {
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

	allCombos := getCombos(pv)
	oldValue := float64(0)

	finalCombos := []combo{}

	for _, c := range allCombos {
		tempCombos := append(finalCombos, c)
		if err = t.Execute(fo, tempCombos); err != nil {
			panic(err)
		}

		updateApproved()
		value := getValue()
		if oldValue < value {
			finalCombos = tempCombos
			oldValue = value
			fmt.Println("--> added combo")
			fmt.Println(c)
			fmt.Println("------------")
		}
		// time.Sleep(2 * time.Second)
		log.Println("----------------------> percent = ", value)
		cleanFile(fo)
		if value == float64(100) {
			break
		}
	}
	if err = t.Execute(fo, finalCombos); err != nil {
		panic(err)
	}

}
