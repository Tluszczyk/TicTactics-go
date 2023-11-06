package log

import (
	golog "log"
)

func Info(message string) {
	golog.Printf("INFO: %s", message)
}

func Warn(message string) {
	golog.Printf("WARNING: %s", message)
}

func Error(message string) {
	golog.Printf("ERROR: %s", message)
}
