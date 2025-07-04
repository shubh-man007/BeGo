// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"
// )

// func worker(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("Process Completed")
// 			return
// 		default:
// 			fmt.Println(time.Now(), "Processing")
// 			time.Sleep(time.Millisecond * 500)
// 		}
// 	}
// }

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go worker(ctx)

// 	time.Sleep(2 * time.Second)

// 	cancel()
// 	fmt.Println("Done")
// }

package main

import (
	"fmt"
	"time"
)

func printLoop() {
	for i := 0; i < 5; i++ {
		fmt.Printf("Processing %v ...\n", i)
		time.Sleep(1 * time.Second) // Simulate work with delay
	}
}

func main() {
	fmt.Println("Starting Goroutine")

	// Goroutines run concurrently and independently from the main function.
	go printLoop()

	// The main goroutine now sleeps for 2 seconds.
	// During this time, the printLoop goroutine runs in the background.
	time.Sleep(2 * time.Second)

	// After 2 seconds, main prints "Done" and exits.
	fmt.Println("Done")

	// IMPORTANT:
	// When main() exits, the entire program ends,
	// even if the printLoop goroutine hasn't finished.
}
