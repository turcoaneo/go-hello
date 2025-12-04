package logger

import (
	"fmt"
	"time"
)

// Log prints a timestamped message with a level
func Log(level string, msg string) {
	fmt.Printf("[%s] %s: %s\n", time.Now().Format("15:04:05"), level, msg)
}

// Info Warn Error Convenience wrappers
func Info(msg string) {
	Log("INFO", msg)
}

//goland:noinspection GoUnusedExportedFunction
func Warn(msg string) {
	Log("WARN", msg)
}

//goland:noinspection GoUnusedExportedFunction
func Error(msg string) {
	Log("ERROR", msg)
}
