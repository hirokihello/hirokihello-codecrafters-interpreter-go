package run

import (
	"fmt"
	"os"
)

func Run() {
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	parser := Parser{
		tokens: tokenize(fileContents),
		index:  0,
	}

	statements := parser.parseStatements()

	env := NewEnv()
	for _, statement := range statements {
		statement.Execute(env)
	}
}
