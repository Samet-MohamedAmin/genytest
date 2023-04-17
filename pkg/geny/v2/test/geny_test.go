package geny

import (
	"reflect"
	"testing"
)

var (
	scoresGiven = []float64{}

	expectedOutput = []float64{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
	}
)

func init() {
	for i := 1; i <= 10; i++ {
		for j := 0; j < 1000; j++ {
			scoresGiven = append(scoresGiven, float64(i))
		}
	}
}

func TestDoTheThing(t *testing.T) {

	scoresInput := make([]float64, len(scoresGiven))
	output := doTheThing(scoresGiven, &scoresInput)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("doTheThing(\n%v\n, scoresInput) = \n%v\n, expected \n%v\n", scoresGiven, output, expectedOutput)
	}

	if !reflect.DeepEqual(scoresGiven, scoresInput) {
		t.Errorf("scoresInput = \n%v\nexpected = %v", scoresInput, scoresGiven)
	}
}

func TestDoTheThing2(t *testing.T) {

	scoresInput := make([]float64, len(scoresGiven))
	output := doTheThing2(scoresGiven, &scoresInput)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("doTheThing(\n%v\n, scoresInput) = \n%v\n, expected \n%v\n", scoresGiven, output, expectedOutput)
	}

	if !reflect.DeepEqual(scoresGiven, scoresInput) {
		t.Errorf("scoresInput = \n%v\nexpected = %v", scoresInput, scoresGiven)
	}
}

func TestDoTheParallelThing(t *testing.T) {

	scoresInput := make([]float64, len(scoresGiven))
	output := doTheThingParallel(scoresGiven, &scoresInput)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("doTheThing(\n%v\n, scoresInput) = \n%v\n, expected \n%v\n", scoresGiven, output, expectedOutput)
	}

	if !reflect.DeepEqual(scoresGiven, scoresInput) {
		t.Errorf("scoresInput = \n%v\nexpected = %v", scoresInput, scoresGiven)
	}
}

func TestDoTheParallelThing2(t *testing.T) {

	scoresInput := make([]float64, len(scoresGiven))
	output := doTheThing2Parallel(scoresGiven, &scoresInput)

	if !reflect.DeepEqual(expectedOutput, output) {
		t.Errorf("doTheThing(\n%v\n, scoresInput) = \n%v\n, expected \n%v\n", scoresGiven, output, expectedOutput)
	}

	if !reflect.DeepEqual(scoresGiven, scoresInput) {
		t.Errorf("scoresInput = \n%v\nexpected = %v", scoresInput, scoresGiven)
	}
}

func BenchmarkTheThing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		scoresInput := make([]float64, len(scoresGiven))
		doTheThing(scoresGiven, &scoresInput)
	}
}

func BenchmarkTheThing2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		scoresInput := make([]float64, len(scoresGiven))
		doTheThing2(scoresGiven, &scoresInput)
	}
}

func BenchmarkTheThingParallel(b *testing.B) {
	// for i := 0; i < b.N; i++ {
	scoresInput := make([]float64, len(scoresGiven))
	doTheThingParallel(scoresGiven, &scoresInput)
	// }
}

func BenchmarkTheThing2Parallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		scoresInput := make([]float64, len(scoresGiven))
		doTheThing2Parallel(scoresGiven, &scoresInput)
	}
}
