package mutex

import (
	"runtime"
	"sync/atomic"
)

type RWMutex struct {
	// m          Mutex
	readState  atomic.Bool
	writeState atomic.Bool
}

func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

func (m *RWMutex) Rlock() {
	// много чтения разрешено пока нет write
	if m.writeState.Load() {
		runtime.Gosched()
	}

}

func (m *RWMutex) RUnlock() {

}

func (m *RWMutex) Wlock() {
	// запрещено пока есть читатели
}

func (m *RWMutex) WUnlock() {

}
