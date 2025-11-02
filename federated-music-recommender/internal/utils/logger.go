package utils

import "fmt"

func Info(msg string) {
	fmt.Printf("🟢 INFO: %s\n", msg)
}

func Warn(msg string) {
	fmt.Printf("🟡 WARN: %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("🔴 ERROR: %s\n", msg)
}
