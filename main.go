package main

import (
	"fmt"
	"runtime"
	"time"
)

var cpuCount = runtime.NumCPU()

// Task type
type Task struct {
	id      int
	payload string
}

// Worker consumes tasks and produces results
func worker(id int, tasks <-chan Task, results chan<- string) {
	for task := range tasks {
		msg := fmt.Sprintf("[Worker %d] processing Task %d (%s) | CPUs: %d | Goroutines: %d",
			id, task.id, task.payload, cpuCount, runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond) // simulate work
		results <- msg
	}
}

func main() {
	fmt.Println("=== Producer–Consumer Cycle ===")
	fmt.Printf("CPUs: %d | Goroutines: %d\n", cpuCount, runtime.NumGoroutine())

	tasks := make(chan Task, 5)     // buffered channel for tasks
	results := make(chan string, 5) // buffered channel for results

	// Launch workers (consumers)
	for i := 1; i <= 3; i++ {
		go worker(i, tasks, results)
	}

	// Producer: send tasks
	for i := 1; i <= 6; i++ {
		tasks <- Task{id: i, payload: fmt.Sprintf("payload-%d", i)}
	}
	close(tasks) // signal no more tasks

	// Collect results
	for i := 1; i <= 6; i++ {
		fmt.Println(<-results)
	}

	fmt.Println("All tasks processed.")
}
