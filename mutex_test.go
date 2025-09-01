package mutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMutex_LockUnlock(t *testing.T) {
	t.Parallel()

	m := NewMutex()

	// Базовый тест lock/unlock
	m.Lock()
	if !m.state.Load() {
		t.Error("Mutex should be locked after Lock()")
	}

	m.Unlock()
	if m.state.Load() {
		t.Error("Mutex should be unlocked after Unlock()")
	}
}

func TestMutex_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	m := NewMutex()
	counter := 0
	var wg sync.WaitGroup
	iterations := 1000

	// Запускаем много горутин для конкурентного доступа
	for range iterations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Lock()
			counter++
			m.Unlock()
		}()
	}

	wg.Wait()

	if counter != iterations {
		t.Errorf("Expected counter %d, got %d", iterations, counter)
	}
}

func TestMutex_TryLockSimulation(t *testing.T) {
	t.Parallel()

	m := NewMutex()

	// Первый lock должен успешно захватиться
	m.Lock()

	// Вторая попытка lock в другой горутине должна блокироваться
	locked := atomic.Bool{}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		m.Lock()
		locked.Store(true)
		m.Unlock()
	}()

	// Даем время на попытку захвата
	time.Sleep(10 * time.Millisecond)
	if locked.Load() {
		t.Error("Second goroutine should be blocked")
	}

	// Разблокируем и проверяем, что вторая горутина смогла захватить
	m.Unlock()
	wg.Wait()

	if !locked.Load() {
		t.Error("Second goroutine should have acquired lock after unlock")
	}
}

func TestMutex_StressTest(t *testing.T) {
	t.Parallel()

	m := NewMutex()
	var wg sync.WaitGroup
	counter := 0
	operations := 100
	task := 10000

	addTask := func() {
		defer wg.Done()
		for range operations {
			m.Lock()
			counter++
			m.Unlock()
		}
	}

	for range task {
		wg.Add(1)
		go addTask()
	}

	wg.Wait()
	expected := task * operations
	if counter != expected {
		t.Errorf("Expected %d, got %d", expected, counter)
	}
}

func TestMutex_UnlockPanic(t *testing.T) {
	t.Parallel()

	m := NewMutex()
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from Unlock without Lock")
			}
		}()
		m.Unlock()
	}()
}
