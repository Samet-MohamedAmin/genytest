package utilities

import (
	"log"
	"sync"
)

type Safeindex struct {
	index int
	mutex sync.Mutex
}

func (sc *Safeindex) Lock() {
	sc.mutex.Lock()
}

func (sc *Safeindex) Unlock() {
	sc.mutex.Unlock()
}

func (sc *Safeindex) SetValue(i int) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.index = i
}

func (sc *Safeindex) Increase() {
	sc.IncreaseBy(1)
	log.Printf("safe value = %d", sc.index)
}

func (sc *Safeindex) IncreaseBy(n int) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.index = sc.index + n
}

func (sc *Safeindex) Value() int {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	return sc.index
}

func (sc *Safeindex) ValueUnsafe() int {
	return sc.index
}
