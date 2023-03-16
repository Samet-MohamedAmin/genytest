package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

const (
	path                  = "../example"
	coverProfile          = "coverage.out"
	htmlPage              = "bla.html"
	possibleValuesFile    = "possible_values.yaml"
	testCasesTemplateFile = "test_cases.yaml.tpl"
	testCasesOutputFile   = path + "/test_cases.yaml"
)

var (
	home, _        = os.LookupEnv("HOME")
	goMutestingBin = home + "/go/bin/go-mutesting"
)

func runCommand(command []string, skipError bool) string {
	cmd := exec.Command(command[0], command[1:]...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(&stdBuffer) //, os.Stdout)

	cmd.Stdout = mw
	// cmd.Stderr = mw

	if err := cmd.Run(); !skipError && err != nil {
		panic(err)
	}

	return stdBuffer.String()
}
