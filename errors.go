package main

import "fmt"

var errorMessages []string = make([]string, 0, 1)

type TrollError struct {
	line int
	message string
}

func (err *TrollError) toString() string {
	return fmt.Sprintf("TrollScript error on line %d. %s", err.line, err.message)
}

func TrollHasError() bool {
	return len(errorMessages) != 0
}

func TrollReportError(line int, message string) {
	errorMessages = append(errorMessages, fmt.Sprintf("TrollScript error on line %d. %s", line, message))
}
