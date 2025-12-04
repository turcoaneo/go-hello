package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Task struct {
	id      int
	payload string
}

func worker(id int, tasks <-chan Task, results chan<- string, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		time.Sleep(500 * time.Millisecond) // simulate heavy load
		msg := fmt.Sprintf("[Worker %d] processed Task %d (%s)", id, task.id, task.payload)
		results <- msg
		logger.Info("Worker finished", "worker", id, "taskID", task.id)
	}
}

// Keyboard input goroutine
func keyboardInput(results chan<- string, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	logger.Info("Keyboard input ready (type something)")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "quit" {
			logger.Info("Keyboard input stopped")
			return
		}
		results <- fmt.Sprintf("[Keyboard] %s", text)
	}
}

// File input goroutine
func fileInput(filename string, results chan<- string, logger *slog.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		logger.Error("Failed to open file", "error", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error("Failed to defer file to closing", "error", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		results <- fmt.Sprintf("[File] %s", line)
		time.Sleep(3000 * time.Millisecond) // simulate delay
	}
	logger.Info("File input finished")
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	tasks := make(chan Task, 5)
	results := make(chan string, 20)

	var wg sync.WaitGroup

	// Workers
	numWorkers := 2
	wg.Add(numWorkers)
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results, logger, &wg)
	}

	// Producer tasks
	for i := 1; i <= 4; i++ {
		tasks <- Task{id: i, payload: fmt.Sprintf("payload-%d", i)}
	}
	close(tasks)

	// Keyboard input
	wg.Add(1)
	go keyboardInput(results, logger, &wg)

	// File input
	wg.Add(1)
	go fileInput("input.txt", results, logger, &wg)

	// Collector goroutine
	go func() {
		wg.Wait()
		close(results)
	}()

	// Consume results
	for msg := range results {
		logger.Info(msg)
	}

	logger.Info("All sources finished.")
}
