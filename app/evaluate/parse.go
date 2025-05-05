package evaluate

import (
	"fmt"
	"os"
)

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	COMMA         = "COMMA"
	DOT           = "DOT"
	MINUS         = "MINUS"
	PLUS          = "PLUS"
	SEMICOLON     = "SEMICOLON"
	STAR          = "STAR"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	SLASH         = "SLASH"
	BANG          = "BANG"
	STRING        = "STRING"
	NUMBER        = "NUMBER"
	NIL           = "NIL"
	TRUE          = "TRUE"
	FALSE         = "FALSE"
	BOOLEAN       = "BOOLEAN"
)

var reservedTokens = map[string]string{
	"(":  LEFT_PAREN,
	")":  RIGHT_PAREN,
	"{":  LEFT_BRACE,
	"}":  RIGHT_BRACE,
	",":  COMMA,
	".":  DOT,
	"-":  MINUS,
	"+":  PLUS,
	";":  SEMICOLON,
	"*":  STAR,
	"=":  EQUAL,
	"==": EQUAL_EQUAL,
	"!=": BANG_EQUAL,
	"<":  LESS,
	"<=": LESS_EQUAL,
	">":  GREATER,
	">=": GREATER_EQUAL,
	"/":  SLASH,
	"!":  BANG,
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
	getValue() EvaluateNode
	getType() string
}

type EvaluateNode struct {
	value     string
	valueType string
}

type StringNode struct {
	Node
	value     string
	tokenType string
}
type NumberNode struct {
	Node
	value     string
	tokenType string
}

type BooleanNode struct {
	Node
	value     string
	tokenType string
}

type NilNode struct {
	Node
	value     string
	tokenType string
}

type Group struct {
	Node
	nodes     []Node
	tokenType string
}

type Unary struct {
	Node
	operator  Token
	right     Node
	tokenType string
}

type Binary struct {
	Node
	left      Node
	operator  Token
	right     Node
	tokenType string
}

func Parse() AST {
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
	return ast
}

// parse して構文木を作成する
func (p *Parser) parse() AST {
	ast := AST{}
	var node Node
	err := error(nil)
	node, _, err = p.parseExpression(0)
	if err != nil {
		os.Exit(65)
	}
	ast.nodes = append(ast.nodes, node)
	return ast
}

func (p *Parser) parseExpression(index int) (Node, int, error) {
	var node Node
	err := error(nil)
	node, index, err = p.parseEquality(index)
	if err != nil {
		return nil, index, err
	}
	return node, index, nil
}

func (p *Parser) parseEquality(index int) (Node, int, error) {
	err := error(nil)
	left, index, err := p.parseComparison(index)

	if err != nil {
		return nil, index, err
	}
	token := p.tokens[index]
	for token.tokenType == EQUAL_EQUAL || token.tokenType == BANG_EQUAL {
		var right Node
		right, index, err = p.parseComparison(index + 1)
		if err != nil {
			return nil, index, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index, nil
}

func (p *Parser) parseComparison(index int) (Node, int, error) {
	err := error(nil)
	var left Node
	left, index, err = p.parseTerm(index)
	if err != nil {
		return nil, index, err
	}
	token := p.tokens[index]
	for token.tokenType == LESS || token.tokenType == LESS_EQUAL ||
		token.tokenType == GREATER || token.tokenType == GREATER_EQUAL {
		var right Node
		right, index, err = p.parseTerm(index + 1)
		if err != nil {
			return nil, index, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}
	return left, index, nil
}

func (p *Parser) parseTerm(index int) (Node, int, error) {
	err := error(nil)
	var left Node
	left, index, err = p.parseFactor(index)

	if err != nil {
		return nil, index, err
	}
	token := p.tokens[index]
	for token.tokenType == PLUS || token.tokenType == MINUS {
		var right Node
		right, index, err = p.parseFactor(index + 1)
		if err != nil {
			return nil, index, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index, nil
}

func (p *Parser) parseFactor(index int) (Node, int, error) {
	err := error(nil)
	var left Node
	left, index, err = p.parseUnary(index)
	if err != nil {
		return nil, index, err
	}

	token := p.tokens[index]
	for token.tokenType == STAR || token.tokenType == SLASH {
		var right Node
		right, index, err = p.parseUnary(index + 1)
		if err != nil {
			return nil, index, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[index]
	}

	return left, index, nil
}

func (p *Parser) parseUnary(index int) (Node, int, error) {
	err := error(nil)
	if index >= len(p.tokens) {
		panic("Index out of range")
	}
	token := p.tokens[index]
	for token.tokenType == BANG || token.tokenType == MINUS {
		var right Node
		right, index, err = p.parseUnary(index + 1)
		if err != nil {
			return nil, index, err
		}
		return &Unary{
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}, index, nil
	}

	return p.parsePrimary(index)
}
func (p *Parser) parsePrimary(index int) (Node, int, error) {
	token := p.tokens[index]

	if token.tokenType == NIL {
		return &NilNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, index + 1, nil
	}

	if token.tokenType == TRUE || token.tokenType == FALSE {
		return &BooleanNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, index + 1, nil
	}

	if token.tokenType == NUMBER {
		return &NumberNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, index + 1, nil
	}

	if token.tokenType == STRING {
		return &StringNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, index + 1, nil
	}

	if token.tokenType == LEFT_PAREN {
		var expression Node
		expression, index, _ = p.parseExpression(index + 1)
		if p.tokens[index].tokenType == "RIGHT_PAREN" {
			return &Group{
				nodes:     []Node{expression},
				tokenType: token.tokenType,
			}, index + 1, nil
		}
		return expression, index + 1, fmt.Errorf("missing right parenthesis")
	}
	return nil, index + 1, fmt.Errorf("unexpected token: %s", token.value)
}
