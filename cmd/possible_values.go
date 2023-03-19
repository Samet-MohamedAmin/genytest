package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const INT_TYPE = "int"

type Combo map[string]string

type rangeValue struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type possibleValue struct {
	Type   string `yaml:"type"`
	Values struct {
		Static []string     `yaml:"static"`
		Ranges []rangeValue `yaml:"ranges"`
	} `yaml:"values"`
}

type valuesByKey map[string]possibleValue

func UnmarshalPossibleValues() valuesByKey {
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

func getAllPossibleValues(value possibleValue) []string {
	s := set{}
	// static
	for _, pv := range value.Values.Static {
		s.add(pv)
	}

	// range
	if value.Type == INT_TYPE {
		for _, r := range value.Values.Ranges {
			for pv := r.Min; pv <= r.Max; pv++ {
				s.add(fmt.Sprint(pv))
			}
		}
	}

	return s.getItems()
}

func cloneCombos(orig *[]Combo) ([]Combo, error) {
	origYAML, err := yaml.Marshal(orig)
	if err != nil {
		return nil, err
	}

	clone := []Combo{}
	if err = yaml.Unmarshal(origYAML, &clone); err != nil {
		return nil, err
	}

	return clone, nil
}

func GetCombos(vmap valuesByKey) (combos []Combo) {
	allValues := map[string][]string{}
	keys := []string{}
	for key, value := range vmap {
		allValues[key] = getAllPossibleValues(value)
		keys = append(keys, key)
	}

	for _, pv := range allValues[keys[0]] {
		combo := map[string]string{
			keys[0]: pv,
		}
		combos = append(combos, combo)
	}

	for _, key := range keys[1:] {
		oldCombos := combos
		combos = []Combo{}
		for _, value := range allValues[key] {
			newCombos, _ := cloneCombos(&oldCombos)
			for _, combo := range newCombos {
				combo[key] = value
			}
			combos = append(combos, newCombos...)
		}
	}

	return
}
