package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to read file %s", err.Error())
	}

	script := string(content)
	scanner := ScannerInit(script)
	scanner.scanTokens()

	if TrollHasError() {
		for msg := range errorMessages {
			log.Print(msg)
		}
		os.Exit(1)
	}

	parser := ParserInit(scanner.tokens)
	expr := parser.parse()

	if TrollHasError() {
		for msg := range errorMessages {
			log.Print(msg)
		}
		os.Exit(1)
	}

	value, trollError := evaluate(expr)
	if trollError != nil {
		log.Print(trollError.toString())
		os.Exit(1)
	}

	log.Print(value.toString())
}
