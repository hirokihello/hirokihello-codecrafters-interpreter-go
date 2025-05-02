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
		errCount := 0;
		for i, x := range fileContents {
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
			} else if x == '-' {
				fmt.Println("MINUS - null")
			} else if x == ';' {
				fmt.Println("SEMICOLON ; null")
			} else if x == '/' {
				fmt.Println("SLASH / null")
			} else if x == '=' {
				if i + 1 < len(fileContents) && fileContents[i+1] == '=' {
					fmt.Println("EQUAL_EQUAL == null")
				} else {
					fmt.Println("EQUAL = null")
				}
			} else {
				fmt.Fprintf(os.Stderr, "[line 1] Error: Unexpected character: %c\n", x)
				errCount++
			}
		}
		fmt.Println("EOF  null")

		if(errCount > 0) {
			os.Exit(65)
		}
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
