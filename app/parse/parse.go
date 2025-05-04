package parse

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"
)

var reservedTokens = map[string]string{
	"(":  "LEFT_PAREN",
	")":  "RIGHT_PAREN",
	"{":  "LEFT_BRACE",
	"}":  "RIGHT_BRACE",
	",":  "COMMA",
	".":  "DOT",
	"-":  "MINUS",
	"+":  "PLUS",
	";":  "SEMICOLON",
	"*":  "STAR",
	"=":  "EQUAL",
	"==": "EQUAL_EQUAL",
	"!=": "BANG_EQUAL",
	"<":  "LESS",
	"<=": "LESS_EQUAL",
	">":  "GREATER",
	">=": "GREATER_EQUAL",
	"/":  "SLASH",
}

var reservedWords = map[string]string{
	"and":    "AND",
	"class":  "CLASS",
	"else":   "ELSE",
	"false":  "FALSE",
	"for":    "FOR",
	"fun":    "FUN",
	"if":     "IF",
	"nil":    "NIL",
	"or":     "OR",
	"print":  "PRINT",
	"return": "RETURN",
	"super":  "SUPER",
	"this":   "THIS",
	"true":   "TRUE",
	"var":    "VAR",
	"while":  "WHILE",
}

type Token struct {
	tokenType string
	value     string
}

// tokenize は、ファイルの内容をトークンに変換します。
// これは、トークンのリストを返します。
func tokenize(fileContents []byte) []Token {
	errCount := 0
	lineCount := 1
	tokens := make([]Token, 0)
	for i := 0; i < len(fileContents); i++ {
		x := fileContents[i]
		if x == '(' || x == ')' || x == '}' || x == '{' || x == '*' || x == '+' || x == '.' || x == ',' ||
			x == '-' || x == ';' {
			tokens = append(tokens, Token{tokenType: reservedTokens[string(x)], value: string(x)})
		} else if x == '=' || x == '!' || x == '<' || x == '>' {
			if i+1 < len(fileContents) && fileContents[i+1] == '=' {
				tokens = append(tokens, Token{
					tokenType: reservedTokens[string(x)+string(fileContents[i+1])],
					value:     string(x) + string(fileContents[i+1]),
				})
				i++
			} else {
				tokens = append(tokens, Token{tokenType: reservedTokens[string(x)], value: string(x)})
			}
		} else if x == '"' {
			string_token := ""

			for i+1 < len(fileContents) && fileContents[i+1] != '"' {
				i++
				string_token += string(fileContents[i])
			}

			if i+1 < len(fileContents) && fileContents[i+1] == '"' {
				tokens = append(tokens,
					Token{
						tokenType: "STRING",
						value:     string_token,
					})
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
			tokens = append(tokens, Token{
				tokenType: "NUMBER",
				value:     number_formatted,
			})
		} else if ('a' <= x && x <= 'z') || x == '_' || ('A' <= x && x <= 'Z') {
			str := ""
			str += string(x)
			for i+1 < len(fileContents) && (('a' <= fileContents[i+1] && fileContents[i+1] <= 'z') ||
				fileContents[i+1] == '_' || ('0' <= fileContents[i+1] && fileContents[i+1] <= '9') || ('A' <= fileContents[i+1] && fileContents[i+1] <= 'Z')) {
				i++
				str += string(fileContents[i])
			}

			if reservedWords[str] != "" {
				tokens = append(tokens,
					Token{
						tokenType: reservedWords[str],
						value:     str,
					})
			} else {
				tokens = append(tokens, Token{
					tokenType: "IDENTIFIER",
					value:     str,
				})
			}
		} else if x == '/' {
			if i+1 < len(fileContents) && fileContents[i+1] == '/' {
				for i+1 < len(fileContents) && fileContents[i+1] != '\n' {
					i++
				}
			} else {
				tokens = append(tokens, Token{
					tokenType: "SLASH",
					value:     string(x),
				})
			}
		} else if x == ' ' || x == '\t' {
			// Ignore whitespace
		} else if x == '\n' {
			lineCount++
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineCount, x)
			errCount++
		}
	}

	tokens = append(tokens, Token{tokenType: "EOF", value: ""})

	// エラーが起こっていた場合は exit code 65 を返す
	if errCount > 0 {
		os.Exit(65)
	}

	return tokens
}

func Parse() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	ast := tokenize(fileContents)

	for _, token := range ast {
		if token.tokenType == "NIL" || token.tokenType == "TRUE" || token.tokenType == "FALSE" {
			fmt.Printf(token.value)
		} else if token.tokenType == "NUMBER" {
			fmt.Printf(token.value)
		} else if token.tokenType == "STRING" {
			fmt.Printf(token.value)
		}
	}
}
