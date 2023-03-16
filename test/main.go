package main

import (
	"genytest/cmd"
)

func main() {

	pv := cmd.GetPossibleValues()
	// fmt.Println(pv)

	// cmd.GenTestCasesFile(pv, cmd.GetTestCoverage)
	cmd.GenTestCasesFile(pv, cmd.GetMutationTestPercent)
}
