package possiblevalues

import (
	"fmt"
	"genytest/config"
	"genytest/utilities"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const INT_TYPE = "int"

type rangeValue struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type PossibleValue struct {
	Type   string `yaml:"type"`
	Values struct {
		Static []string     `yaml:"static"`
		Ranges []rangeValue `yaml:"ranges"`
	} `yaml:"values"`
}

func (v PossibleValue) GetAllPossibleValues() []string {
	s := utilities.Set{}
	// static
	for _, pv := range v.Values.Static {
		s.Add(pv)
	}

	// range
	if v.Type == INT_TYPE {
		for _, r := range v.Values.Ranges {
			for pv := r.Min; pv <= r.Max; pv++ {
				s.Add(fmt.Sprint(pv))
			}
		}
	}

	return s.GetItems()
}

func unmarshalPossibleValues() (valuesMap map[string]PossibleValue) {
	filename, _ := filepath.Abs(config.PossibleValuesFile)
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &valuesMap)
	if err != nil {
		panic(err)
	}

	return valuesMap
}

func GetAllValues() (allValues map[string][]string, keys []string) {
	vmap := unmarshalPossibleValues()

	for key, value := range vmap {
		allValues[key] = value.GetAllPossibleValues()
		keys = append(keys, key)
	}
	return
}
