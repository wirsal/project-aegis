package logger

import (
	"fmt"
	"log"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m" // merah (error)
	colorGreen  = "\033[92m" // hijau terang (info)
	colorYellow = "\033[93m" // kuning terang (warn)
	colorCyan   = "\033[36m" // cyan (debug)
)

type logLevel struct {
	consolePrefix string
	filePrefix    string
	color         string
}

// Mapping log level
var levels = map[string]logLevel{
	"INFO":  {"[INFO]", "INFO", colorGreen},
	"ERROR": {"[ERROR]", "ERROR", colorRed},
	"DEBUG": {"[DEBUG]", "DEBUG", colorCyan},
	"WARN":  {"[WARN]", "WARN", colorYellow},
}

// Public logging functions
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLogAll("INFO", msg, nil)
}

func Error(msg string, err error) {
	writeLogAll("ERROR", msg, err)
}

func Debug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLogAll("DEBUG", msg, nil)
}

func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLogAll("WARN", msg, nil)
}

// Core logging
func writeLogAll(levelKey, msg string, err error) {
	level, ok := levels[levelKey]
	if !ok {
		level = logLevel{"[UNKNOWN]", levelKey, ""}
	}

	logToConsole(level, msg, err)
	logToFile(level, msg, err)
}

func logToConsole(level logLevel, msg string, err error) {
	if err != nil {
		log.Printf(level.color+"%s %s: %v"+colorReset, level.consolePrefix, msg, err)
	} else {
		log.Printf(level.color+"%s %s"+colorReset, level.consolePrefix, msg)
	}
}

func logToFile(level logLevel, msg string, err error) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	var formatted string
	if err != nil {
		formatted = fmt.Sprintf("%s [%s] %s: %v\n", timestamp, level.filePrefix, msg, err)
	} else {
		formatted = fmt.Sprintf("%s [%s] %s\n", timestamp, level.filePrefix, msg)
	}

	if _, writeErr := writeLog(formatted); writeErr != nil {
		log.Printf("[ERROR] Failed to write log to file: %v\n", writeErr)
	}
}

func writeLog(msg string) (bool, error) {
	// logPath := "" //filepath.Join(GetCurrentLoggingDir(), viper.GetString("file.Log.Logging"))

	// if err := os.MkdirAll(GetCurrentLoggingDir(), os.ModePerm); err != nil {
	// 	return false, fmt.Errorf("failed to create directory: %w", err)
	// }

	// file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to open file: %w", err)
	// }
	// defer file.Close()

	// if _, err := file.WriteString(msg); err != nil {
	// 	return false, fmt.Errorf("failed to write data: %w", err)
	// }

	// if err := file.Sync(); err != nil {
	// 	return false, fmt.Errorf("failed to sync file: %w", err)
	// }

	return true, nil
}
