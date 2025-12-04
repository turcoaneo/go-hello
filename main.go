package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"
)

var cpuCount = runtime.NumCPU()

type Task struct {
	id      int
	payload string
}

func worker(id int, tasks <-chan Task, results chan<- string, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done() // signal completion when worker exits

	for task := range tasks {
		msg := fmt.Sprintf("[Worker %d] processing Task %d (%s) | CPUs: %d | Goroutines: %d",
			id, task.id, task.payload, cpuCount, runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond) // simulate work
		results <- msg
		logger.Info("Task finished",
			"worker", id,
			"taskID", task.id,
			"payload", task.payload,
		)
	}
}

func main() {
	// Choose handler (JSON or text)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info("Producer–Consumer Cycle Started",
		"cpus", cpuCount,
		"goroutines", runtime.NumGoroutine(),
	)

	tasks := make(chan Task, 6)
	results := make(chan string, 6)

	var wg sync.WaitGroup

	// Launch workers
	numWorkers := 3
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results, logger, &wg)
	}

	// Producer: send tasks
	for i := 1; i <= 6; i++ {
		tasks <- Task{id: i, payload: fmt.Sprintf("payload-%d", i)}
		logger.Info("Produced task", "taskID", i)
	}
	close(tasks) // signal no more tasks

	// Wait for workers to finish
	wg.Wait()
	close(results)

	// Collect results
	for result := range results {
		logger.Info(result)
	}

	logger.Info("All tasks processed.")
}
