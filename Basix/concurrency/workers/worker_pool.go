package main

import (
	"fmt"
	"sync"
)

// chan<- int is send-only: the function can only send values, not receive or close the channel,
// which makes intent explicit and prevents accidental misuse.

// jobs <-chan int restricts this worker to receiving tasks only, making channel ownership explicit.

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing task %d\n", id, job)
		results <- 2 * job
		fmt.Printf("Worker %d completed task %d\n", id, job)
	}
}

func main() {
	var wg sync.WaitGroup
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	for i := range 3 {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	for j := range 5 {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)

	for k := range results {
		fmt.Printf("Result: %d\n", k)
	}

	fmt.Println("Tasks Processed !")

}
