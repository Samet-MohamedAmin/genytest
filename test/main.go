package main

import (
	"fmt"
	"genytest/cmd"
)

func main() {

	pv := cmd.UnmarshalPossibleValues()

	allCombos := cmd.GetCombos(pv)

	combos := []cmd.Combo{}

	fmt.Println("---> Coverage")
	combos = cmd.GenTestCasesFile(allCombos, combos, cmd.GetTestCoverage)
	fmt.Println("---> Mutation")
	cmd.GenTestCasesFile(allCombos, combos, cmd.GetMutationTestPercent)
}
