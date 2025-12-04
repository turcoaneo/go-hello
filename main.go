package main

import (
	"fmt"
	"go-hello/logger"
	"runtime"
	"time"
)

var cpuCount = runtime.NumCPU()

type Task struct {
	id      int
	payload string
}

func worker(id int, tasks <-chan Task, results chan<- string) {
	for task := range tasks {
		msg := fmt.Sprintf("[Worker %d] processing Task %d (%s) | CPUs: %d | Goroutines: %d",
			id, task.id, task.payload, cpuCount, runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond)
		results <- msg
		logger.Info(fmt.Sprintf("Worker %d finished Task %d", id, task.id))
	}
}

func main() {
	logger.Info("=== Producer–Consumer Cycle Started ===")
	logger.Info(fmt.Sprintf("CPUs: %d | Goroutines: %d", cpuCount, runtime.NumGoroutine()))

	tasks := make(chan Task, 5)
	results := make(chan string, 5)

	for i := 1; i <= 3; i++ {
		go worker(i, tasks, results)
	}

	for i := 1; i <= 6; i++ {
		tasks <- Task{id: i, payload: fmt.Sprintf("payload-%d", i)}
		logger.Info(fmt.Sprintf("Produced Task %d", i))
	}
	close(tasks)

	for i := 1; i <= 6; i++ {
		logger.Info(<-results)
	}

	logger.Info("All tasks processed.")
}
