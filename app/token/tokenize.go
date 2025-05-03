package token

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"
)

func Tokenize() {
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
				number_formatted := ""
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
					tmp_num, _ := strconv.ParseFloat(number_token, 64)
					number_formatted = strconv.FormatFloat(tmp_num, 'g', -1, 64)
					if math.Mod(tmp_num, 1) == 0 {
						number_formatted = strconv.FormatFloat(tmp_num, 'g', -1, 64) + ".0"
					}

				} else {
					number_formatted = number_token + ".0"
				}
				fmt.Printf("NUMBER %s %s\n", number_token, number_formatted)
			} else if ('a' <= x && x <= 'z') || x == '_' || ('A' <= x && x <= 'Z') {
				str := ""
				str += string(x)
				for i+1 < len(fileContents) && (('a' <= fileContents[i+1] && fileContents[i+1] <= 'z') ||
					fileContents[i+1] == '_' || ('0' <= fileContents[i+1] && fileContents[i+1] <= '9') || ('A' <= fileContents[i+1] && fileContents[i+1] <= 'Z')) {
					i++
					str += string(fileContents[i])
				}

				if str == "and" {
					fmt.Println("AND and null")
				} else if str == "class" {
					fmt.Println("CLASS class null")
				} else if str == "else" {
					fmt.Println("ELSE else null")
				} else if str == "false" {
					fmt.Println("FALSE false null")
				} else if str == "for" {
					fmt.Println("FOR for null")
				} else if str == "fun" {
					fmt.Println("FUN fun null")
				} else if str == "if" {
					fmt.Println("IF if null")
				} else if str == "nil" {
					fmt.Println("NIL nil null")
				} else if str == "or" {
					fmt.Println("OR or null")
				} else if str == "print" {
					fmt.Println("PRINT print null")
				} else if str == "return" {
					fmt.Println("RETURN return null")
				} else if str == "super" {
					fmt.Println("SUPER super null")
				} else if str == "this" {
					fmt.Println("THIS this null")
				} else if str == "true" {
					fmt.Println("TRUE true null")
				} else if str == "var" {
					fmt.Println("VAR var null")
				} else if str == "while" {
					fmt.Println("WHILE while null")
				} else {
					fmt.Printf("IDENTIFIER %s null\n", str)
				}
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
