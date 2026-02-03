package main

import (
	"fmt"
	"sync"
	"time"
)

// chan error        // bidirectional channel
// chan<- error      // send-only channel
// chan<- error is send-only: the function can only send values, not receive or close the channel,
// which makes intent explicit and prevents accidental misuse.

func workSim(id int, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	if id%2 == 0 {
		errChan <- fmt.Errorf("Error processing task %d\n", id)
		return
	}

	time.Sleep(time.Second)
	fmt.Printf("Processed task %d\n", id)
}

func main() {
	var wg sync.WaitGroup
	chanErr := make(chan error, 4)

	for i := range 4 {
		wg.Add(1)
		go workSim(i, &wg, chanErr)
	}

	wg.Wait()
	close(chanErr)

	for err := range chanErr {
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}

	// An open channel with no senders causes receivers to block forever.
	// close(chanErr)

	fmt.Println("Processed Tasks !")
}
