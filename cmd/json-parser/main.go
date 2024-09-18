package main

import (
	"fmt"
	"json-parser/pkg/parser"
	"json-parser/pkg/tokenizer"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	var t *tokenizer.Tokenizer
	switch len(args) {
	case 1:
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Unable to read file: %v", err)
		}
		defer file.Close()
		t = tokenizer.NewTokenizerFromReader(file)
	default:
		if isInputFromPipe() {
			t = tokenizer.NewTokenizerFromReader(os.Stdin)
		} else {
			printUsageAndExit()
		}
	}
	if parser.Parse(t) {
		os.Exit(0)
	} else {
		fmt.Println("Invalid json.")
		os.Exit(1)
	}
}

func isInputFromPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func printUsageAndExit() {
	fmt.Println("Usage: json-parser <filename> or cat <filename> | json-parser")
	os.Exit(1)
}
