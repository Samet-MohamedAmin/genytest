package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	coverageRegex = "coverage: ([\\d.]+)% of statements"
)

func GetTestCoverage() float64 {
	filename, err := filepath.Abs(path)

	if err != nil {
		fmt.Println(err)
	}

	command := []string{
		"go",
		"test",
		"-coverprofile=" + coverProfile,
		"-covermode=count",
		filename,
	}

	tcResult := runCommand(command, false)

	re := regexp.MustCompile(coverageRegex)

	parts := re.FindStringSubmatch(tcResult)

	var coveragePercent float64

	if coveragePercent, err = strconv.ParseFloat(parts[1], 64); err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(coveragePercent)

	command = []string{
		"go",
		"tool",
		"cover",
		"-html=" + coverProfile,
		"-o=bla.html",
	}
	runCommand(command, false)

	return coveragePercent
}
