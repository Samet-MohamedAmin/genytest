package scores

import (
	"genytest/config"
	"genytest/utilities"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	coverageRegex = "coverage: ([\\d.]+)% of statements"
)

var (
	coverageArgs = []string{
		"go",
		"test",
		"-coverprofile=" + config.CoverProfile,
		"-covermode=count",
	}
)

type Coverage struct{}

func (Coverage) GetValue(subpath string) (coveragePercent float64) {

	args := append(coverageArgs, subpath+"/"+config.Path)
	c := exec.Command(args[0], args[1:]...)
	c.Dir = subpath

	tcResult := utilities.RunCommand(c, false)

	re := regexp.MustCompile(coverageRegex)

	parts := re.FindStringSubmatch(tcResult)

	var err error

	if coveragePercent, err = strconv.ParseFloat(parts[1], 64); err != nil {
		log.Fatalf(err.Error())
	}

	return
}
