package main

import (
	"fmt"
	"time"
)

func worker(id int, work chan int, done chan bool) {
	for i := range work {
		fmt.Printf("Worker %d processing task ID: %d\n", id, i)
		time.Sleep(time.Second)
		done <- true
	}
}

func main() {
	ch := make(chan string)

	go func() {
		ch <- "Sending message via channel"
	}()

	msg := <-ch
	fmt.Println(msg)

	fmt.Println("----")

	// buffered channels
	ch1 := make(chan int, 2)
	ch1 <- 1
	ch1 <- 2

	fmt.Println(<-ch1)
	fmt.Println(<-ch1)

	fmt.Println("----")

	// closing channels
	ch2 := make(chan int)

	go func() {
		for i := range 5 {
			ch2 <- i
		}
		close(ch2)
	}()

	for v := range ch2 {
		fmt.Println(v)
	}

	fmt.Println("----")

	// select statements
	// select statements allow you to handle multiple channel operations simultaneously.
	// They enable a goroutine to wait on multiple communication operations. The select statement blocks until one of its cases can run, then it executes that case.
	// If multiple cases can proceed, one is chosen at random.
	// Default case can be used to have a non-blocking select
	ch3 := make(chan string)
	ch4 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- "Sent into ch3"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch4 <- "Sent into ch4"
	}()

	for range 2 {
		select {
		case msg := <-ch3:
			fmt.Printf("ch3: %s\n", msg)
		case msg := <-ch4:
			fmt.Printf("ch4: %s\n", msg)
		}
	}

	fmt.Println("----")

	// task parallelization
	work := make(chan int, 5)
	done := make(chan bool, 5)

	for i := range 5 {
		go worker(i, work, done)
	}

	for j := range 5 {
		work <- j
	}

	close(work)

	for range 5 {
		<-done
	}

	fmt.Println("Processes consumed !")

}
