package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

const (
	path               = "/var/home/mohamedamin/Desktop/test/GildedRose-Refactoring-Kata/go/gildedrose"
	coverProfile       = "coverage.out"
	htmlPage           = "bla.html"
	possibleValuesFile = "possible_values.yaml"
)

var (
	home, _        = os.LookupEnv("HOME")
	goMutestingBin = home + "/go/bin/go-mutesting"
)

type set map[string]any

func runCommand(command []string, skipError bool) string {
	cmd := exec.Command(command[0], command[1:]...)

	cmd.Dir = path

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(&stdBuffer) //, os.Stdout)

	cmd.Stdout = mw
	// cmd.Stderr = mw

	if err := cmd.Run(); !skipError && err != nil {
		panic(err)
	}

	return stdBuffer.String()
}

func (s set) add(item string) {
	s[item] = nil
}

func (s set) getItems() (items []string) {
	for key := range s {
		items = append(items, key)
	}
	return
}
