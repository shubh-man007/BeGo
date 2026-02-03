package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var val int32 = 0
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&val, 1)
		}()
	}

	wg.Wait()

	fmt.Printf("Value: %d\n", val)

	swapped := atomic.CompareAndSwapInt32(&val, 1000, 1001)
	fmt.Printf("Swapped: %v, Value: %d\n", swapped, val)
}
