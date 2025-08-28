package main

import (
	"fmt"
	"sync"
	"time"

	mutex "github.com/Ontair/mutex/internal"
)

const task = 10000



func testRWMutex() {
	fmt.Println("Testing RWMutex...")
	m := mutex.NewRWMutex()
	data := make(map[int]string)
	var wg sync.WaitGroup
	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.RLock()
			fmt.Printf("Reader %d: read %d items\n", id, len(data))
			m.RUnlock()
		}(i)
	}
	// Writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.Lock()
			data[id] = fmt.Sprintf("value-%d", id)
			fmt.Printf("Writer %d: wrote data\n", id)
			m.Unlock()
		}(i)
	}

	// Readers
	for i := 11; i < 40; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.RLock()
			fmt.Printf("Reader %d: read %d items\n", id, len(data))
			m.RUnlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final data size: %d\n", len(data))
}

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


	time.Sleep(100 * time.Millisecond)
	testRWMutex()
}
