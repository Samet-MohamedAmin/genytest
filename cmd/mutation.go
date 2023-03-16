package cmd

import (
	"log"
	"regexp"
	"strconv"
)

const (
	mutationRegex = "The mutation score is (\\S+)"
)

func GetMutationTestPercent() float64 {
	command := []string{
		goMutestingBin,
		path,
	}

	tcResult := runCommand(command, true)

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
