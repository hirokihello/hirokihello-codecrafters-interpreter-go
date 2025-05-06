package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/evaluate"
	"github.com/codecrafters-io/interpreter-starter-go/app/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/run"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "parse" && command != "tokenize" && command != "evaluate" && command != "run" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	if command == "tokenize" {
		token.Tokenize()
	}

	if command == "parse" {
		parse.Parse()
	}

	if command == "evaluate" {
		evaluate.Evaluate()
	}

	if command == "run" {
		run.Run()
	}
}
