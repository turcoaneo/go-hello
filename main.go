package main

import (
	"fmt"
	"runtime"
	"time"
)

// Global variable holding CPU count
var cpuCount = runtime.NumCPU()

func worker(id int) {
	fmt.Printf("[Worker %d] CPUs: %d | Goroutines: %d\n",
		id, cpuCount, runtime.NumGoroutine())

	for i := 0; i < 3; i++ {
		fmt.Printf("[Worker %d] iteration %d\n", id, i)
		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	fmt.Println("=== System Info at Start ===")
	fmt.Printf("CPUs: %d | Goroutines: %d\n", cpuCount, runtime.NumGoroutine())

	// Launch workers
	for i := 1; i <= 3; i++ {
		go worker(i)
	}

	// Print after spawning
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("After spawning workers: CPUs: %d | Goroutines: %d\n",
		cpuCount, runtime.NumGoroutine())

	// Keep main alive
	time.Sleep(2 * time.Second)
	fmt.Println("All goroutines finished.")
}
