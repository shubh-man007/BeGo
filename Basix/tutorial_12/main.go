package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string) // Unbuffered channel

	go func() {
		fmt.Println("Within the func goroutine")
		time.Sleep(1 * time.Second)
		ch <- "Sending data through channel" // Blocks until main receives
	}()

	msg := <-ch // Blocks until goroutine sends
	fmt.Println(msg)

	ch1 := make(chan int, 3) // Buffered channel with capacity 3

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i // May block if buffer is full
		}
		close(ch1) // Required to end range loop
	}()

	for v := range ch1 { // Receives until channel is closed
		fmt.Println(v)
	}
}
