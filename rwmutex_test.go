package mutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRWMutex_BasicReadWrite(t *testing.T) {
	m := NewRWMutex()

	// Чтение должно работать когда нет писателей
	m.RLock()
	if m.readCount.Load() != 1 {
		t.Error("Should have 1 reader")
	}
	m.RUnlock()

	// Письмо должно блокировать чтение
	m.Lock()
	if !m.writeState.Load() {
		t.Error("Should be in write state")
	}
	m.Unlock()
}

func TestRWMutex_MultipleReaders(t *testing.T) {
	m := NewRWMutex()
	var wg sync.WaitGroup
	readers := 10

	// Запускаем несколько читателей одновременно
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			m.RLock()
			if m.readCount.Load() < 1 {
				t.Errorf("Reader %d: should have readers count >= 1", id)
			}
			time.Sleep(1 * time.Millisecond)
			m.RUnlock()
		}(i)
	}

	wg.Wait()

	// После всех читателей счетчик должен быть 0
	if m.readCount.Load() != 0 {
		t.Error("All readers should have finished")
	}
}

func TestRWMutex_WriterBlocksReaders(t *testing.T) {
	m := NewRWMutex()

	// Писатель захватывает мьютекс
	m.Lock()

	// Читатель должен заблокироваться
	readerActive := atomic.Bool{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		m.RLock()
		readerActive.Store(true)
		m.RUnlock()
	}()

	// Даем время читателю попытаться
	time.Sleep(20 * time.Millisecond)
	if readerActive.Load() {
		t.Error("Reader should be blocked while writer is active")
	}

	// Разблокируем писателя
	m.Unlock()

	// Теперь читатель должен разблокироваться
	wg.Wait()
	if !readerActive.Load() {
		t.Error("Reader should have acquired lock after writer")
	}
}

func TestRWMutex_ReadersBlockWriter(t *testing.T) {
	m := NewRWMutex()

	// Читатель захватывает мьютекс
	m.RLock()

	// Писатель должен заблокироваться
	writerActive := atomic.Bool{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		m.Lock()
		writerActive.Store(true)
		m.Unlock()
	}()

	// Даем время писателю попытаться
	time.Sleep(20 * time.Millisecond)
	if writerActive.Load() {
		t.Error("Writer should be blocked while readers are active")
	}

	// Разблокируем читателя
	m.RUnlock()

	// Теперь писатель должен разблокироваться
	wg.Wait()
	if !writerActive.Load() {
		t.Error("Writer should have acquired lock after readers")
	}
}

func TestRWMutex_StressTest(t *testing.T) {
	m := NewRWMutex()
	data := make(map[int]int)
	var readCount, writeCount atomic.Int64
	var wg sync.WaitGroup

	// Запускаем читателей
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				m.RLock()
				_ = len(data) // Читаем данные
				readCount.Add(1)
				m.RUnlock()
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Запускаем писателей
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				m.Lock()
				data[id*100+j] = j
				writeCount.Add(1)
				m.Unlock()
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()

	if readCount.Load() != 5000 {
		t.Errorf("Expected 5000 reads, got %d", readCount.Load())
	}
	if writeCount.Load() != 200 {
		t.Errorf("Expected 200 writes, got %d", writeCount.Load())
	}
}

func TestRMMutex_RUnlockPanic(t *testing.T) {
	rm := NewRWMutex()

	go func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from RUnlock without Lock")
			}
		}()
		rm.RUnlock()
	}()
}
