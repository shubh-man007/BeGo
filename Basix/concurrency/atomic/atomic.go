package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type User struct {
	Name string
	Age  int
}

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

	// Store user struct
	var value atomic.Value
	user := User{Name: "Little John", Age: 30}
	value.Store(user)

	// Load and print.
	current := value.Load().(User)
	fmt.Printf("Current user: %+v\n", current)

	// Update with new user data.
	user2 := User{Name: "Vera", Age: 36}
	value.Store(user2)

	// Load and print updated data.
	updated := value.Load().(User)
	fmt.Printf("Updated user: %+v\n", updated)

}
