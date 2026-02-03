package main

import (
	"context"
	"fmt"
	"time"
)

func worker1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done")
			return
		default:
			fmt.Println("Worker working ...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func worker2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Timed Out")
			return

		default:
			fmt.Println("Worker on Time ...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func worker3(ctx context.Context) {
	if taskID, ok := ctx.Value("taskID").(int); ok {
		fmt.Printf("Processing task: %d\n", taskID)
	} else {
		fmt.Println("No task ID found in context")
	}
}

func main() {
	ctx1, cancel := context.WithCancel(context.Background())
	go worker1(ctx1)
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(time.Second) // let worker sleep !

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go worker2(ctx2)
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(time.Second)

	ctx3 := context.WithValue(context.Background(), "taskID", 3)
	worker3(ctx3)
}
