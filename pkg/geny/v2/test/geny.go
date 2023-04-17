package geny

import (
	"genytest/config"
	"genytest/pkg/combo"
	"genytest/scores"
	"html/template"
	"os"
	"path"
	"sync"
	"time"
)

func writeTemplate(t *template.Template, pathFile string, data any) {

	// open output file
	fo, err := os.Create(pathFile)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if err := t.Execute(fo, data); err != nil {
		panic(err)
	}
}

func getScore(scoresGiven []float64, index int) float64 {
	time.Sleep(time.Millisecond * 50)
	return scoresGiven[index]
}

func getScore2(scoresGiven []float64, scoresInput *[]float64, index int) float64 {
	if index < len(*scoresInput) && (*scoresInput)[index] != 0 {
		return (*scoresInput)[index]
	}
	return getScore(scoresGiven, index)
}

func shrinkScores(scoresInput *[]float64) (output []float64) {

	prevScore := float64(0)

	for _, s := range *scoresInput {
		if prevScore < s {
			output = append(output, s)
		}
		prevScore = s
	}

	return
}

func doTheSubThing(scoresGiven []float64, scoresInput *[]float64, firstIndex int, lastIndex int) {
	if lastIndex-firstIndex < 0 {
		return
	}

	if firstIndex == lastIndex {
		score := getScore(scoresGiven, firstIndex)
		(*scoresInput)[firstIndex] = score
		return
	}

	scoreFirst := getScore(scoresGiven, firstIndex)
	scoreLast := getScore(scoresGiven, lastIndex)

	if scoreFirst == scoreLast {
		for i := firstIndex; i <= lastIndex; i++ {
			(*scoresInput)[i] = scoreFirst
		}
		return
	}

	medium := (firstIndex + lastIndex) / 2
	if medium < firstIndex || medium > lastIndex {
		return
	}
	doTheSubThing(scoresGiven, scoresInput, firstIndex+1, medium)
	doTheSubThing(scoresGiven, scoresInput, medium+1, lastIndex-1)
}

func doTheThing(scoresGiven []float64, scoresInput *[]float64) []float64 {
	if len(*scoresInput) < 2 {
		return nil
	}

	doTheSubThing(scoresGiven, scoresInput, 0, len(*scoresInput)-1)

	return shrinkScores(scoresInput)
}

func doTheSubThing2(scoresGiven []float64, scoresInput *[]float64, firstIndex int, lastIndex int) {
	const MAX = 10

	if lastIndex-firstIndex < 0 {
		return
	}

	i := firstIndex
	for ; i <= lastIndex; i++ {
		score := getScore(scoresGiven, i)
		(*scoresInput)[i] = score
		if score == MAX {
			break
		}
	}

	for ; i <= lastIndex; i++ {
		(*scoresInput)[i] = MAX
	}

}

func doTheThing2(scoresGiven []float64, scoresInput *[]float64) []float64 {
	if len(*scoresInput) < 2 {
		return nil
	}

	doTheSubThing2(scoresGiven, scoresInput, 0, len(*scoresInput)-1)

	return shrinkScores(scoresInput)
}

func setScore(scoresGiven []float64, scoresInput *[]float64, index int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	if index < len(*scoresInput) && (*scoresInput)[index] != 0 {
		return
	}
	(*scoresInput)[index] = getScore(scoresGiven, index)
}

func doTheSubThingParallel(scoresGiven []float64, scoresInput *[]float64, firstIndex int, lastIndex int, wg *sync.WaitGroup) {
	defer wg.Done()

	if lastIndex-firstIndex < 0 {
		return
	}

	if firstIndex == lastIndex {
		setScore(scoresGiven, scoresInput, firstIndex, nil)
		return
	}

	var wg2 sync.WaitGroup
	wg2.Add(2)
	setScore(scoresGiven, scoresInput, firstIndex, &wg2)
	setScore(scoresGiven, scoresInput, lastIndex, &wg2)
	wg2.Wait()

	if (*scoresInput)[firstIndex] == (*scoresInput)[lastIndex] {
		for i := firstIndex; i <= lastIndex; i++ {
			(*scoresInput)[i] = (*scoresInput)[lastIndex]
		}
		return
	}

	medium := (firstIndex + lastIndex) / 2

	wg.Add(2)
	go doTheSubThingParallel(scoresGiven, scoresInput, firstIndex, medium, wg)
	go doTheSubThingParallel(scoresGiven, scoresInput, medium+1, lastIndex, wg)

}

func doTheThingParallel(scoresGiven []float64, scoresInput *[]float64) []float64 {
	if len(*scoresInput) < 2 {
		return nil
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go doTheSubThingParallel(scoresGiven, scoresInput, 0, len(*scoresInput)-1, &wg)

	wg.Wait()

	return shrinkScores(scoresInput)
}

func doTheSubThing2Parallel(scoresGiven []float64, scoresInput *[]float64, firstIndex int, lastIndex int) {

	const MAX = 10

	if lastIndex-firstIndex < 0 {
		return
	}
	var wg sync.WaitGroup

	mutex := sync.Mutex{}
	maxIndex := 0

	for i := 0; i <= lastIndex; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			mutex.Lock()
			if (*scoresInput)[maxIndex] == MAX && maxIndex >= index {
				(*scoresInput)[index] = MAX
			}
			mutex.Unlock()

			score := getScore(scoresGiven, index)
			(*scoresInput)[index] = score

			mutex.Lock()
			maxValue := (*scoresInput)[maxIndex]
			if score > maxValue {
				maxIndex = index
			}
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
}

func doTheThing2Parallel(scoresGiven []float64, scoresInput *[]float64) []float64 {
	if len(*scoresInput) < 2 {
		return nil
	}

	doTheSubThing2Parallel(scoresGiven, scoresInput, 0, len(*scoresInput)-1)

	return shrinkScores(scoresInput)
}
func GenTestCasesFile(allCombos combo.Combos, finalCombos combo.Combos, customTest scores.CustomTest) combo.Combos {
	if len(allCombos) < 2 {
		return allCombos
	}

	t, err := template.ParseFiles(config.TestCasesTemplateFile)
	if err != nil {
		panic(err)
	}

	outputPath := path.Join(config.BasePath, config.Path, config.TestCasesOutputFile)
	writeTemplate(t, outputPath, finalCombos)

	return finalCombos
}
