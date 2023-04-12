package main

import (
	"fmt"
	"genytest/pkg/combo"
	"genytest/pkg/geny"
	"genytest/scores"
)

func main() {

	allCombos := combo.Combos{}.GenCombos()

	combos := combo.Combos{}

	// fmt.Println("---> Coverage")
	// coverage := geny.Coverage{}
	// combos = geny.GenTestCasesFile(allCombos, combos, coverage)

	fmt.Println("---> Mutation")
	mutation := scores.Mutation{}
	geny.GenTestCasesFile(allCombos, combos, mutation)
}
