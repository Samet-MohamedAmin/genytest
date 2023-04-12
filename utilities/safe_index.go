package utilities

import "sync"

type Safeindex struct {
	index int
	mutex sync.Mutex
}

func (sc *Safeindex) SetValue(i int) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.index = i
}

func (sc *Safeindex) Value() int {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	return sc.index
}
