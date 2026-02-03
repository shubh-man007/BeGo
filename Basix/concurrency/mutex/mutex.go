package main

import (
	"fmt"
	"sync"
)

var value int = 0

func Increment(wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	for range 1000 {
		mu.Lock()
		value++
		mu.Unlock()
	}
}

func IncrementReadRW(wg *sync.WaitGroup, mu *sync.RWMutex) {
	defer wg.Done()

	mu.RLock()
	v := value
	mu.RUnlock()

	fmt.Println("Value:", v)
}

func IncrementWriteRW(wg *sync.WaitGroup, mu *sync.RWMutex) {
	defer wg.Done()

	for range 1000 {
		mu.Lock()
		value++
		mu.Unlock()
	}
}

func main() {
	// Mutex for synchroization
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go Increment(&wg, &mu)
	}

	wg.Wait()
	fmt.Printf("Updated value (Mutex): %d\n", value)

	// Mutex for read-heavy operations:
	var wg2 sync.WaitGroup
	var muRW sync.RWMutex

	wg2.Add(3)
	go IncrementReadRW(&wg2, &muRW)
	go IncrementWriteRW(&wg2, &muRW)
	go IncrementReadRW(&wg2, &muRW)

	wg2.Wait()
	fmt.Printf("Updated value (RWMutex): %d\n", value)
}
