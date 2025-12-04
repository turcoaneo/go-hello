package main

import (
	"fmt"
	"runtime"
	"time"
)

var cpuCount = runtime.NumCPU()

func worker(id int, ch chan string) {
	for i := 0; i < 3; i++ {
		msg := fmt.Sprintf("[Worker %d] CPUs: %d | Goroutines: %d | iteration %d",
			id, cpuCount, runtime.NumGoroutine(), i)
		ch <- msg // send message into channel
		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== System Info at Start ===")
	fmt.Printf("CPUs: %d | Goroutines: %d\n", cpuCount, runtime.NumGoroutine())

	ch := make(chan string)

	// Launch workers
	for i := 1; i <= 3; i++ {
		go worker(i, ch)
	}

	// Collect results
	for i := 0; i < 9; i++ { // 3 iterations × 3 workers
		fmt.Println(<-ch) // receive from channel
	}

	fmt.Println("All goroutines finished.")
}
