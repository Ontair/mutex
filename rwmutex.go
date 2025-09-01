package mutex

import (
	"runtime"
	"sync/atomic"
)

type RWMutex struct {
	readCount  atomic.Int64
	writeState atomic.Bool
}

func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

func (m *RWMutex) RLock() {
	// много чтения разрешено пока нет write
	for {
		if m.writeState.Load() {
			runtime.Gosched()
			continue
		}

		m.readCount.Add(1)
		// Вторая проверка, на тот случай, если после первой проверки и до добавления читателя не встроился писатель
		if !m.writeState.Load() {
			return // Не встроился
		}

		// пистель встролся - откатываемся
		m.readCount.Add(-1)
	}
}

func (m *RWMutex) RUnlock() {
	if m.readCount.Add(-1) < 0 {
		panic("RUnlock of unlocked RWMutex")
	}
}

func (m *RWMutex) Lock() {
	// запрещено пока есть читатели
	retries := 0
	for !m.writeState.CompareAndSwap(unlocked, locked) {
		retries++
		if retries > fastCheckNumber {
			runtime.Gosched() // даём попытку другим горутинам поработать
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
