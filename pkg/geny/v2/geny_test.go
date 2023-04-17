package geny_test

import (
	"fmt"
	"genytest/pkg/combo"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

var (
	filename = "medium_test_cases.yaml"
)

type testCase struct {
	Name    string `json:"name" yaml:"name"`
	SellIn  int    `json:"sellIn" yaml:"sellIn"`
	Quality int    `json:"quality" yaml:"quality"`
}

func readTestCases() (testCases []testCase) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &testCases)
	if err != nil {
		panic(err)
	}

	return
}

func TestSequentialFilter(t *testing.T) {

	testCases := readTestCases()
	combos := combo.Combos{}
	for _, tc := range testCases {
		combo := combo.Combo{}
		combo.Items["name"] = tc.Name
		combo.Items["sellIn"] = fmt.Sprint(tc.SellIn)
		combo.Items["quality"] = fmt.Sprint(tc.Quality)
		combos = append(combos, &combo)
	}

	fmt.Println(combos)

}
