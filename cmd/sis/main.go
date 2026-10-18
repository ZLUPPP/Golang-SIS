package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"log.analizer/internal/application"
	"log.analizer/internal/domain"
	"log.analizer/internal/parser"
	"log.analizer/internal/storage"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Log Analyzer SIS
Usage:
  sis -p <file> [-f <format>]

Options:
`)

		flag.PrintDefaults()

		io.WriteString(flag.CommandLine.Output(), `
Supported format tokens:
  %Y        year
  %m        month
  %d        day
  %H        hour
  %M        minute
  %S        second
  %LEVEL    log level: INFO, WARN or ERROR
  %MESSAGE  log message

Examples:
  sis -p "test.log"
  sis -p "test.log" -f "%Y-%m-%d [%LEVEL] %MESSAGE"
  sis -p "test.log" -f "%Y-%m-%d %H:%M:%S [%LEVEL] %MESSAGE"
`)
	}

	path := flag.String("p", "", "path to the log file")
	format := flag.String(
		"f",
		"%Y-%m-%d %LEVEL %MESSAGE",
		"format of the log file",
	)

	flag.Parse()

	if *path == "" {
		fmt.Fprintln(os.Stderr, "file path is required: use -p")
		os.Exit(1)
	}

	entries, err := readLogFile(*path, *format)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read logs:", err)
		os.Exit(1)
	}

	repository := storage.NewMemoryRepository()
	service := application.NewService(repository)

	task, err := service.CreateTask(entries)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create task:", err)
		os.Exit(1)
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "process task:", err)
		os.Exit(1)
	}

	result := processed.Result

	fmt.Printf("Task ID: %d\n", processed.ID)
	fmt.Printf("Status: %s\n", processed.Status)
	fmt.Printf("Total: %d\n", result.Total)
	fmt.Printf("INFO: %d\n", result.InfoCount)
	fmt.Printf("WARN: %d\n", result.WarnCount)
	fmt.Printf("ERROR: %d\n", result.ErrorCount)
	fmt.Printf("Most common error: %s\n", result.MostCommonError)
}

func readLogFile(path string, format string) ([]domain.LogEntry, error) {
	lineParser, err := parser.New(format)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []domain.LogEntry

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := lineParser.ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, domain.ErrEmptyInput
	}

	return entries, nil
}
