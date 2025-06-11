package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	pontoFilePathEnv = os.Getenv("PONTO_FILE_PATH")
	showExitHour     = os.Getenv("PONTO_SHOW_EXITH")
	pontoLogPathEnv  = os.Getenv("PONTO_LOG_PATH")
	logger           *log.Logger
)

const (
	targetHours       = 8
	timeFormat        = "15:04"
	pontoFileBaseName = ".ponto"
	logFileBaseName   = ".ponto.log"
	prefixWorked      = "w-"
	prefixRemaining   = "r-"
)

func main() {
	setupLogger()

	var pontoFileName string
	if pontoFilePathEnv != "" {
		pontoFileName = filepath.Join(pontoFilePathEnv, pontoFileBaseName)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logFatal("Cannot determine user home dir for ponto file: %v", err)
		}
		pontoFileName = filepath.Join(homeDir, pontoFileBaseName)
	}

	file, err := os.Open(pontoFileName)
	if err != nil {
		if os.IsNotExist(err) {
			logFatal("Point file '%s' not found.", pontoFileName)
		} else {
			logFatal("Error opening file '%s': %v", pontoFileName, err)
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totalWorkedDuration := calcTotalHoursWorked(scanner)

	if err := scanner.Err(); err != nil {
		logFatal("Error reading file: %v", err)
	}

	if totalWorkedDuration == 0 {
		logWarning("No valid time entries found in ponto file.")
	}

	targetDuration := time.Duration(targetHours) * time.Hour

	var outputMinutes int
	var prefix, currentStatusMsg string

	if totalWorkedDuration >= targetDuration {
		outputMinutes = int(totalWorkedDuration.Minutes())
		prefix = prefixWorked
	} else {
		prefix = prefixRemaining
		remainingDuration := targetDuration - totalWorkedDuration
		outputMinutes = int(remainingDuration.Minutes())
	}

	hours := outputMinutes / 60
	minutes := outputMinutes % 60
	currentStatusMsg = fmt.Sprintf("%s%02d:%02d", prefix, hours, minutes)
	fmt.Println(showExitHour)
	if showExitHour == "true" {
		currentStatusMsg = buildExitHourMessage(totalWorkedDuration, targetDuration, currentStatusMsg)
	}

	fmt.Println(currentStatusMsg)
}

func setupLogger() {
	var logFilePath string

	if pontoLogPathEnv != "" {
		logFilePath = filepath.Join(pontoLogPathEnv, logFileBaseName)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Failed to determine user home directory:", err)
			os.Exit(1)
		}
		logFilePath = filepath.Join(homeDir, logFileBaseName)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}

	logger = log.New(logFile, "PONTO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func logWarning(format string, a ...interface{}) {
	if logger != nil {
		logger.Printf("WARN: "+format, a...)
	}
}

func logFatal(format string, a ...interface{}) {
	if logger != nil {
		logger.Printf("FATAL: "+format, a...)
	}
	os.Exit(1)
}

func buildExitHourMessage(totalWorkedDuration, targetDuration time.Duration, baseStatusMsg string) string {
	now := time.Now()

	if totalWorkedDuration >= targetDuration {
		return fmt.Sprintf("%s (Goal Met)", baseStatusMsg)
	}

	remainingTimeNeeded := targetDuration - totalWorkedDuration
	projectedExitTime := now.Add(remainingTimeNeeded)

	if projectedExitTime.Before(now) {
		return fmt.Sprintf("%s (Time Missed)", baseStatusMsg)
	}

	return fmt.Sprintf("%s (Exit: %s)", baseStatusMsg, projectedExitTime.Format(timeFormat))
}

func parseTime(baseDate time.Time, timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time '%s': %w", timeStr, err)
	}
	return time.Date(
		baseDate.Year(),
		baseDate.Month(),
		baseDate.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0,
		0,
		baseDate.Location(),
	), nil
}

func calcTotalHoursWorked(scanner *bufio.Scanner) time.Duration {
	totalWorkedDuration := time.Duration(0)
	now := time.Now()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		numParts := len(parts)

		switch numParts {
		case 1:
			entryStr := parts[0]
			entryTimeCombinedToday, err := parseTime(now, entryStr)
			if err != nil {
				logWarning("Invalid time format in line '%s'. Ignoring: %v", line, err)
				continue
			}
			finalEntryTime := entryTimeCombinedToday
			if finalEntryTime.After(now) {
				finalEntryTime = finalEntryTime.Add(-24 * time.Hour)
			}
			duration := now.Sub(finalEntryTime)
			if duration > 12*time.Hour {
				logWarning("Unrealistic session duration (%v) in line: '%s'. Ignoring.", duration, line)
				continue
			}
			totalWorkedDuration += duration

		case 2:
			entryStr := parts[0]
			exitStr := parts[1]

			entryTime, err := parseTime(now, entryStr)
			if err != nil {
				logWarning("Invalid entry time in line '%s'. Ignoring: %v", line, err)
				continue
			}
			exitTime, err := parseTime(now, exitStr)
			if err != nil {
				logWarning("Invalid exit time in line '%s'. Ignoring: %v", line, err)
				continue
			}
			if exitTime.Before(entryTime) {
				exitTime = exitTime.Add(24 * time.Hour)
			}
			totalWorkedDuration += exitTime.Sub(entryTime)

		default:
			logWarning("Invalid line format '%s'. Expected 'HH:MM' or 'HH:MM HH:MM'. Ignoring.", line)
		}
	}

	return totalWorkedDuration
}
