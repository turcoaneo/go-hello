package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

var cpuCount = runtime.NumCPU()

type Task struct {
	id      int
	payload string
}

func worker(id int, tasks <-chan Task, results chan<- string, logger *slog.Logger) {
	for task := range tasks {
		msg := fmt.Sprintf("[Worker %d] processing Task %d (%s) | CPUs: %d | Goroutines: %d",
			id, task.id, task.payload, cpuCount, runtime.NumGoroutine())
		time.Sleep(500 * time.Millisecond)
		results <- msg
		logger.Info("Task finished",
			"worker", id,
			"taskID", task.id,
			"payload", task.payload,
		)
	}
}

func main() {
	// Choose handler based on config
	var logger *slog.Logger
	if //goland:noinspection GoBoolExpressions
	useJSONLogging {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	logger.Info("Producer–Consumer Cycle Started",
		"cpus", cpuCount,
		"goroutines", runtime.NumGoroutine(),
	)

	tasks := make(chan Task, 5)
	results := make(chan string, 5)

	for i := 1; i <= 3; i++ {
		go worker(i, tasks, results, logger)
	}

	for i := 1; i <= 6; i++ {
		tasks <- Task{id: i, payload: fmt.Sprintf("payload-%d", i)}
		logger.Info("Produced task", "taskID", i)
	}
	close(tasks)

	for i := 1; i <= 6; i++ {
		logger.Info(<-results)
	}

	logger.Info("All tasks processed.")
}
