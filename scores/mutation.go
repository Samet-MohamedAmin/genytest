package scores

import (
	"genytest/config"
	"genytest/utilities"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	mutationRegex = "The mutation score is (\\S+)"
)

var (
	home, _        = os.LookupEnv("HOME")
	goMutestingBin = home + "/go/bin/bin/go-mutesting"
	mutationArgs   = []string{
		goMutestingBin,
		config.Path,
	}
)

type Mutation struct {
}

func (Mutation) GetValue(subpath string) float64 {

	args := append(mutationArgs, subpath+"/"+config.Path)

	c := exec.Command(args[0], args[1:]...)
	c.Dir = subpath

	tcResult := utilities.RunCommand(c, false)

	re := regexp.MustCompile(mutationRegex)

	parts := re.FindStringSubmatch(tcResult)

	var mutationScore float64
	var err error

	if mutationScore, err = strconv.ParseFloat(parts[1], 64); err != nil {
		log.Fatalf(err.Error())
	}

	// fmt.Printf("Mutation Score = %f\n", mutationScore)

	return mutationScore * 100
}
