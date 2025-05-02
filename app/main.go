package main

import (
	"fmt"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		for _, x := range fileContents {
			if x == '(' {
				fmt.Println("LEFT_PAREN ( null")
			} else if x == ')' {
				fmt.Println("RIGHT_PAREN ) null")
			} else if x == '}' {
				fmt.Println("RIGHT_BRACE } null")
			} else if x == '{' {
				fmt.Println("LEFT_BRACE { null")
			} else if x == '*' {
				fmt.Println("STAR * null")
			} else if x == '+' {
				fmt.Println("PLUS + null")
			} else if x == '.' {
				fmt.Println("DOT . null")
			} else if x == ',' {
				fmt.Println("COMMA , null")
			}
		}
		fmt.Println("EOF  null")
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
