package example_test

import (
	"genytest/example"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
	"gopkg.in/yaml.v3"
)

type testTemplate struct {
	Id     string `yaml:"id"`
	Input  int    `yaml:"input"`
	Output int    `Yaml:"output"`
}

const testCasesFile = "test_cases.yaml"

func readTestCases() []testTemplate {
	filename, _ := filepath.Abs(testCasesFile)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var testCases []testTemplate

	err = yaml.Unmarshal(yamlFile, &testCases)
	if err != nil {
		panic(err)
	}

	return testCases
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

func Transform(testCases any) string {
	tc := testCases.(testTemplate)

	tc.Output = example.Bla(tc.Input)

	output, err := yaml.Marshal([]testTemplate{tc})
	if err != nil {
		panic(err)
	}

	return string(output)
}

func TestBla(t *testing.T) {

	testCases := readTestCases()

	approvals.VerifyAll(t, funcName(), testCases, Transform)

}
