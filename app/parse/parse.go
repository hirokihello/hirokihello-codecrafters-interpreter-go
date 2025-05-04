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

type BooleanNode struct {
	Node
	value string
}

type NilNode struct {
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

type Binary struct {
	Node
	left     Node
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
	var node Node
	node, _ = p.parseExpression(0)
	ast.nodes = append(ast.nodes, node)
	return ast
}

func (p *Parser) parseExpression(index int) (Node, int) {
	var node Node
	node, index = p.parseEquality(index)
	return node, index
}

func (p *Parser) parseEquality(index int) (Node, int) {
	left, index := p.parseComparison(index)

	token := p.tokens[index]
	for token.tokenType == "EQUAL_EQUAL" || token.tokenType == "BANG_EQUAL" {
		var right Node
		right, index = p.parseComparison(index + 1)
		left = &Binary{
			left:     left,
			operator: token,
			right:    right,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index
}

func (p *Parser) parseComparison(index int) (Node, int) {
	var left Node
	left, index = p.parseTerm(index)

	token := p.tokens[index]
	for token.tokenType == "LESS" || token.tokenType == "LESS_EQUAL" ||
		token.tokenType == "GREATER" || token.tokenType == "GREATER_EQUAL" {
		var right Node
		right, index = p.parseTerm(index + 1)
		left = &Binary{
			left:     left,
			operator: token,
			right:    right,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}
	return left, index
}

func (p *Parser) parseTerm(index int) (Node, int) {
	var left Node
	left, index = p.parseFactor(index)

	token := p.tokens[index]
	for token.tokenType == "PLUS" || token.tokenType == "MINUS" {
		var right Node
		right, index = p.parseFactor(index + 1)
		left = &Binary{
			left:     left,
			operator: token,
			right:    right,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index
}

func (p *Parser) parseFactor(index int) (Node, int) {
	var left Node
	left, index = p.parseUnary(index)

	token := p.tokens[index]
	for token.tokenType == "STAR" || token.tokenType == "SLASH" {
		var right Node
		right, index = p.parseUnary(index + 1)
		left = &Binary{
			left:     left,
			operator: token,
			right:    right,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index
}

func (p *Parser) parseUnary(index int) (Node, int) {
	if index >= len(p.tokens) {
		panic ("Index out of range")
	}
	token := p.tokens[index]
	for token.tokenType == "BANG" || token.tokenType == "MINUS" {
		var right Node
		right, index = p.parseUnary(index + 1)
		return &Unary{
			operator: token,
			right:    right,
		}, index
	}

	return p.parsePrimary(index)
}
func (p *Parser) parsePrimary(index int) (Node, int) {
	token := p.tokens[index]

	if token.tokenType == "NIL" {
		return &NilNode{
			value: token.value,
		}, index + 1
	}

	if token.tokenType == "TRUE" || token.tokenType == "FALSE" {
		return &BooleanNode{
			value: token.value,
		}, index + 1
	}

	if token.tokenType == "NUMBER" {
		return &NumberNode{
			value: token.value,
		}, index + 1
	}

	if token.tokenType == "STRING" {
		return &StringNode{
			value: token.value,
		}, index + 1
	}

	if token.tokenType == "LEFT_PAREN" {
		var expression Node
		expression, index = p.parseExpression(index + 1)
		if p.tokens[index].tokenType == "RIGHT_PAREN" {
			return &Group{
				nodes: []Node{expression},
			}, index + 1
		}
		return expression, index + 1
	}

	fmt.Printf("Unexpected token: %s\n", token.value)
	fmt.Printf("Unexpected tokenType: %s\n", token.tokenType)
	panic("Unhandled primary expression case")
}

// Print methods for AST and nodes
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
func (n *BooleanNode) Print() {
	io.WriteString(os.Stdout, n.value)
}
func (n *NilNode) Print() {
	io.WriteString(os.Stdout, n.value)
}

func (b *Binary) Print() {
	io.WriteString(os.Stdout, "(")
	io.WriteString(os.Stdout, b.operator.value)
	io.WriteString(os.Stdout, " ")
	b.left.Print()
	io.WriteString(os.Stdout, " ")
	b.right.Print()
	io.WriteString(os.Stdout, ")")
}
