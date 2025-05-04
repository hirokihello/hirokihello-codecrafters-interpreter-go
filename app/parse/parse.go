package parse

import (
	"fmt"
	"io"
	"os"
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
	"!":  "BANG",
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
type Parser struct {
	tokens []Token
	index  int
}

type AST struct {
	nodes []Node
}

type Node interface {
	Print()
}

type StringNode struct {
	Node
	value string
}
type NumberNode struct {
	Node
	value string
}

type Group struct {
	Node
	nodes []Node
}

type Unary struct {
	Node
	operator Token
	right    Node
}

func Parse() {
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

	ast := parser.parse()
	ast.Print()
}

// parse して構文木を作成する
func (p *Parser) parse() AST {
	ast := AST{}
	var n []Node
	for i := 0; i < len(p.tokens); i++ {
		i, n = tokenToNode(p.tokens, i)
		ast.nodes = append(ast.nodes, n...)
	}

	return ast
}

func tokenToNode(tokens []Token, i int) (int, []Node) {
	token := tokens[i]
	nodes := make([]Node, 0)

	switch token.tokenType {
	case "EOF":
		return i, nodes
	case "STRING":
		nodes = append(nodes, &StringNode{value: token.value})
	case "NUMBER":
		nodes = append(nodes, &NumberNode{value: token.value})
	case "LEFT_PAREN":
		group := &Group{}
		for i+1 < len(tokens) {
			pt := tokens[i+1]
			if pt.tokenType == "RIGHT_PAREN" {
				i++
				break
			} else {
				var n []Node
				// 再帰で探す
				// こいつを使用することで、group の nodes に RIGHT_PAREN までのノードを追加する
				i, n = tokenToNode(tokens, i+1)
				group.nodes = append(group.nodes, n...)
			}
		}
		if i >= len(tokens) {
			fmt.Errorf("invalid input. missing )")
			os.Exit(65)
			// TODO: error
		}
		nodes = append(nodes, group)
	case "BANG":
		u := &Unary{}
		u.operator = token
		var n []Node
		i, n = tokenToNode(tokens, i+1)
		u.right = n[0]
		nodes = append(nodes, u)
	case "MINUS":
		u := &Unary{}
		u.operator = token
		var n []Node
		i, n = tokenToNode(tokens, i+1)
		if 1 < len(n)  {
			fmt.Fprintf(os.Stderr, "invalid input. too many nodes")
			os.Exit(65)
		}
		u.right = n[0]
		nodes = append(nodes, u)
	default:
		nodes = append(nodes, &StringNode{value: token.value})
		// fmt.Fprintf(os.Stderr, "Unknown token: %s\n", token.tokenType)
	}
	return i, nodes
}

func (a *AST) Print() {
	for _, n := range a.nodes {
		n.Print()
		io.WriteString(os.Stdout, "\n")
	}
}

func (u *Unary) Print() {
	io.WriteString(os.Stdout, "(")
	io.WriteString(os.Stdout, u.operator.value)
	io.WriteString(os.Stdout, " ")
	u.right.Print()
	io.WriteString(os.Stdout, ")")
}

func (g *Group) Print() {
	io.WriteString(os.Stdout, "(")
	io.WriteString(os.Stdout, "group ")
	for _, n := range g.nodes {
		n.Print()
	}
	io.WriteString(os.Stdout, ")")
}
func (s *StringNode) Print() {
	io.WriteString(os.Stdout, s.value)
}
func (n *NumberNode) Print() {
	io.WriteString(os.Stdout, n.value)
}