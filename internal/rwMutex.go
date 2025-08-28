package mutex

import (
	"runtime"
	// "sync"
	"sync/atomic"
)

type RWMutex struct {

	// m          Mutex
	readCount  atomic.Int64
	writeState atomic.Bool
}

func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

func (m *RWMutex) RLock() {
	for {
		// много чтения разрешено пока нет write
		if m.writeState.Load() {
			runtime.Gosched()
			continue
		}

		m.readCount.Add(1)
		// Вторая проверка (двойная проверка)
		if !m.writeState.Load() {
			return // Успех!
		}

		// Если писатель появился между первой и второй проверкой
		m.readCount.Add(-1) // Откатываемся
	}
}

func (m *RWMutex) RUnlock() {
	m.readCount.Add(-1)
}

func (m *RWMutex) Lock() {
	// запрещено пока есть читатели
	retries := 0
	for !m.writeState.CompareAndSwap(unlocked, locked) {
		retries++
		if retries > fastCheckNumber {
			runtime.Gosched() // даём попытку другим горутинам
			retries = 0
		}
	}

	for m.readCount.Load() > 0 {
		runtime.Gosched()
	}

}

func (m *RWMutex) Unlock() {
	m.writeState.Store(unlocked)
}
