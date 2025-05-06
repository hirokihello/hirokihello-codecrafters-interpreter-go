package run

import (
	"fmt"
	"os"
	"strconv"
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
	PRINT         = "PRINT"
	EOF           = "EOF"
	VAR           = "VAR"
	IDENTIFIER    = "IDENTIFIER"
)

var reservedTokens = map[string]string{
	"(":     LEFT_PAREN,
	")":     RIGHT_PAREN,
	"{":     LEFT_BRACE,
	"}":     RIGHT_BRACE,
	",":     COMMA,
	".":     DOT,
	"-":     MINUS,
	"+":     PLUS,
	";":     SEMICOLON,
	"*":     STAR,
	"=":     EQUAL,
	"==":    EQUAL_EQUAL,
	"!=":    BANG_EQUAL,
	"<":     LESS,
	"<=":    LESS_EQUAL,
	">":     GREATER,
	">=":    GREATER_EQUAL,
	"/":     SLASH,
	"!":     BANG,
	"print": PRINT,
	"var":   VAR,
}

type Token struct {
	tokenType string
	value     string
}
type Parser struct {
	tokens []Token
	index  int
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

type IdentifierNode struct {
	Node
	value     string
	tokenType string
}

// parse して構文木を作成する
func (p *Parser) parseStatements() []Statement {
	statements := make([]Statement, 0)
	for p.index < len(p.tokens) && p.tokens[p.index].tokenType != EOF {
		statement := p.parseStatement()
		if statement == nil {
			fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
			panic("error while parsing")
		}
		statements = append(statements, statement)
	}
	return statements
}

func (p *Parser) parseStatement() Statement {
	if (p.tokens[p.index].tokenType == PRINT) && (p.index+2 < len(p.tokens)) {
		p.index++
		expr, err := p.parseExpression()
		if p.tokens[p.index].value == ";" {
			p.index++
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(65)
			}
			return &PrintStatement{
				expr: expr,
			}
		}

		panic("; is missing")
	} else if p.tokens[p.index].tokenType == VAR && (p.index+3 < len(p.tokens)) {
		p.index++
		varName := p.tokens[p.index].value
		p.index++
		if p.tokens[p.index].value != "=" {
			panic("= is missing")
		}
		p.index++
		varValue, err := p.parseExpression()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		if p.tokens[p.index].value != ";" {
			panic("; is missing")
		}
		p.index++

		return &VariableStatement{
			expr:    varValue,
			varName: varName,
		}
	}

	// ただの式。特に何かをしているわけではない。
	expression, err := p.parseExpression()
	if p.tokens[p.index].value == ";" {
		p.index++
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		return &ExpressionStatement{
			expr: expression,
		}
	}

	panic("; is missing")

}

func (p *Parser) parseExpression() (Node, error) {
	var node Node
	err := error(nil)
	node, err = p.parseEquality()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *Parser) parseEquality() (Node, error) {
	err := error(nil)
	left, err := p.parseComparison()

	if err != nil {
		return nil, err
	}
	token := p.tokens[p.index]
	for token.tokenType == EQUAL_EQUAL || token.tokenType == BANG_EQUAL {
		var right Node
		p.index++
		right, err = p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[p.index]
	}

	return left, nil
}

func (p *Parser) parseComparison() (Node, error) {
	err := error(nil)
	var left Node
	left, err = p.parseTerm()
	if err != nil {
		return nil, err
	}
	token := p.tokens[p.index]
	for token.tokenType == LESS || token.tokenType == LESS_EQUAL ||
		token.tokenType == GREATER || token.tokenType == GREATER_EQUAL {
		var right Node
		p.index++
		right, err = p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[p.index]
	}
	return left, nil
}

func (p *Parser) parseTerm() (Node, error) {
	err := error(nil)
	var left Node
	left, err = p.parseFactor()

	if err != nil {
		return nil, err
	}
	token := p.tokens[p.index]
	for token.tokenType == PLUS || token.tokenType == MINUS {
		var right Node
		p.index++
		right, err = p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[p.index]
	}

	return left, nil
}

func (p *Parser) parseFactor() (Node, error) {
	err := error(nil)
	var left Node
	left, err = p.parseUnary()
	if err != nil {
		return nil, err
	}

	token := p.tokens[p.index]
	for token.tokenType == STAR || token.tokenType == SLASH {
		var right Node
		p.index++
		right, err = p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &Binary{
			left:      left,
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}
		// 次のループに備えて token を更新する
		token = p.tokens[p.index]
	}

	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	err := error(nil)
	if p.index >= len(p.tokens) {
		panic("Index out of range")
	}
	token := p.tokens[p.index]
	for token.tokenType == BANG || token.tokenType == MINUS {
		var right Node
		p.index++
		right, err = p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &Unary{
			operator:  token,
			right:     right,
			tokenType: token.tokenType,
		}, nil
	}

	return p.parsePrimary()
}
func (p *Parser) parsePrimary() (Node, error) {
	token := p.tokens[p.index]

	if token.tokenType == NIL {
		p.index++
		return &NilNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, nil
	}

	if token.tokenType == TRUE || token.tokenType == FALSE {
		p.index++
		return &BooleanNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, nil
	}

	if token.tokenType == NUMBER {
		p.index++
		return &NumberNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, nil
	}

	if token.tokenType == STRING {
		p.index++
		return &StringNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, nil
	}

	if token.tokenType == LEFT_PAREN {
		var expression Node
		p.index++
		expression, _ = p.parseExpression()
		if p.tokens[p.index].tokenType == "RIGHT_PAREN" {
			p.index++
			return &Group{
				nodes:     []Node{expression},
				tokenType: token.tokenType,
			}, nil
		}
		return expression, fmt.Errorf("missing right parenthesis")
	}

	if token.tokenType == IDENTIFIER {
		p.index++
		return &IdentifierNode{
			value:     token.value,
			tokenType: token.tokenType,
		}, nil
	}

	return &NilNode{}, fmt.Errorf("unknown expression")
}

func (u *Unary) getValue() EvaluateNode {
	if u.operator.tokenType == MINUS {
		if u.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operand must be a number.")
			os.Exit(70)
		}
		num, _ := strconv.Atoi(u.right.getValue().value)
		return EvaluateNode{
			value:     strconv.FormatFloat(float64(num*-1), 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if u.operator.tokenType == BANG {
		if u.right.getValue().value == "false" || u.right.getValue().value == "" || u.right.getValue().value == "nil" {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	}

	panic("Unknown operator: " + u.operator.tokenType)
}

func (g *Group) getValue() EvaluateNode {
	if len(g.nodes) == 1 {
		return EvaluateNode{
			value:     g.nodes[0].getValue().value,
			valueType: g.nodes[0].getValue().valueType,
		}
	} else {
		values := ""
		for i, n := range g.nodes {
			if i != 0 {
				values += " "
			}
			values += n.getValue().value
		}
		return EvaluateNode{
			value:     values,
			valueType: STRING,
		}
	}
}

func (i *IdentifierNode) getValue() EvaluateNode {
	variables := getGlobalEnv()

	if _, ok := variables.variables[i.value]; !ok {
		fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n", i.value)
		os.Exit(70)
	}

	return EvaluateNode{
		value:     variables.variables[i.value].value,
		valueType: variables.variables[i.value].valueType,
	}
}

func (s *StringNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     s.value,
		valueType: STRING,
	}
}
func (n *NumberNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NUMBER,
	}
}
func (b *BooleanNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     b.value,
		valueType: BOOLEAN,
	}
}
func (n *NilNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NIL,
	}
}

func (b *Binary) getValue() EvaluateNode {
	if b.operator.tokenType == PLUS {
		if b.left.getValue().valueType != b.right.getValue().valueType {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if b.left.getValue().valueType == STRING {
			return EvaluateNode{
				value:     b.left.getValue().value + b.right.getValue().value,
				valueType: STRING,
			}
		} else if b.left.getValue().valueType == NUMBER {
			left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
			right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
			return EvaluateNode{
				value:     strconv.FormatFloat(left+right, 'f', -1, 64),
				valueType: NUMBER,
			}
		}
	}

	left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
	right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
	if b.operator.tokenType == SLASH {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left/right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == STAR {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left*right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == MINUS {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left-right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == GREATER {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if left > right {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	} else if b.operator.tokenType == GREATER_EQUAL {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if left >= right {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	} else if b.operator.tokenType == LESS {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if left < right {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	} else if b.operator.tokenType == LESS_EQUAL {
		if b.left.getValue().valueType != NUMBER || b.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if left <= right {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	} else if b.operator.tokenType == EQUAL_EQUAL {
		if b.left.getValue().value == b.right.getValue().value && b.left.getValue().valueType == b.right.getValue().valueType {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	} else if b.operator.tokenType == BANG_EQUAL {
		if b.left.getValue().value != b.right.getValue().value || b.left.getValue().valueType != b.right.getValue().valueType {
			return EvaluateNode{
				value:     "true",
				valueType: BOOLEAN,
			}
		} else {
			return EvaluateNode{
				value:     "false",
				valueType: BOOLEAN,
			}
		}
	}

	panic("Unknown operator: " + b.operator.tokenType)
}
