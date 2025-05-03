package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
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
		errCount := 0
		lineCount := 1
		for i := 0; i < len(fileContents); i++ {
			x := fileContents[i]
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
				if i+1 < len(fileContents) && fileContents[i+1] == '/' {
					for i+1 < len(fileContents) && fileContents[i+1] != '\n' {
						i++
					}
				} else {
					fmt.Println("SLASH / null")
				}
			} else if x == '=' {
				if i+1 < len(fileContents) && fileContents[i+1] == '=' {
					fmt.Println("EQUAL_EQUAL == null")
					i++
				} else {
					fmt.Println("EQUAL = null")
				}
			} else if x == '!' {
				if i+1 < len(fileContents) && fileContents[i+1] == '=' {
					fmt.Println("BANG_EQUAL != null")
					i++
				} else {
					fmt.Println("BANG ! null")
				}
			} else if x == '<' {
				if i+1 < len(fileContents) && fileContents[i+1] == '=' {
					fmt.Println("LESS_EQUAL <= null")
					i++
				} else {
					fmt.Println("LESS < null")
				}
			} else if x == '>' {
				if i+1 < len(fileContents) && fileContents[i+1] == '=' {
					fmt.Println("GREATER_EQUAL >= null")
					i++
				} else {
					fmt.Println("GREATER > null")
				}
			} else if x == ' ' || x == '\t' {
				// Ignore whitespace
			} else if x == '\n' {
				lineCount++
			} else if x == '"' {
				string_token := ""

				for i+1 < len(fileContents) && fileContents[i+1] != '"' {
					i++
					string_token += string(fileContents[i])
				}

				if i+1 < len(fileContents) && fileContents[i+1] == '"' {
					fmt.Printf("STRING \"%s\" %s\n", string_token, string_token)
					i++
				} else if i+1 == len(fileContents) {
					errCount++
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", lineCount)
				}
			} else if unicode.IsDigit(rune(x)) {
				number_token := string(x)
				num := 0.0
				formatted := ""
				num_answer := ""
				for i+1 < len(fileContents) && unicode.IsDigit(rune(fileContents[i+1])) {
					i++
					number_token += string(fileContents[i])
				}
				if i+1 < len(fileContents) && fileContents[i+1] == '.' {
					i++
					number_token += string(fileContents[i])
					for i+1 < len(fileContents) && unicode.IsDigit(rune(fileContents[i+1])) {
						i++
						number_token += string(fileContents[i])
					}

					num, _ := strconv.ParseFloat(number_token, 64)
					formatted = strconv.FormatFloat(num, 'g', -1, 64)
					num_answer = formatted
				} else {
					number_token += ".0"
					num, _ = strconv.ParseFloat(number_token, 64)
					num_answer = number_token
				}
				fmt.Printf("NUMBER %f %s\n", num, num_answer)

			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineCount, x)
				errCount++
			}
		}
		fmt.Println("EOF  null")

		if errCount > 0 {
			os.Exit(65)
		}
	} else {
		fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	}
}
