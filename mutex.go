package mutex

import (
	"runtime"
	"sync/atomic"
)

const (
	locked   = true
	unlocked = false
)

var fastCheckNumber = 3 // быстрая попытка захватить память

type Mutex struct {
	state atomic.Bool
}

func NewMutex() *Mutex {
	return &Mutex{}
}

func (m *Mutex) Lock() {
	retries := 0
	for !m.state.CompareAndSwap(unlocked, locked) {
		retries++
		if retries > fastCheckNumber {
			runtime.Gosched() // даём попытку другим горутинам
			retries = 0
		}
	}
}

func (m *Mutex) Unlock() {
	if !m.state.CompareAndSwap(locked, unlocked) {
		panic("unlock of unlocked mutex")
	}
}
