package main

import (
	"fmt"
	"sync"

	mutex "github.com/Ontair/mutex/internal"
)

const task = 10000

func main() {
	

	m := mutex.NewMutex()
	var wg sync.WaitGroup
	counter := 0

	wg.Add(task)

	addTask := func() {
		defer wg.Done()
		m.Lock()
		counter++
		m.Unlock()
	}

	for i := 0; i < task; i++ {
		go addTask()
	}

	wg.Wait()
	fmt.Println(counter)
}
