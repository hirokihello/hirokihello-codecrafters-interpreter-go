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
	ASSIGNMENT    = "ASSIGNMENT"
	IF            = "IF"
	ELSE          = "ELSE"
	OR            = "OR"
	AND           = "AND"
	WHILE         = "WHILE"
	FOR           = "FOR"
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
	"if":    IF,
	"else":  ELSE,
	"or":    OR,
	"and":   AND,
	"while": WHILE,
	"for":   FOR,
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
	getValue(env *Env) EvaluateNode
	getType() string
}

type EvaluateNode struct {
	value     string
	valueType string
}

type AssignmentNode struct {
	Node
	varName   string
	value     Node
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
	if p.index >= len(p.tokens) {
		panic("Index out of range")
	} else if p.tokens[p.index].tokenType == IF {
		p.index++
		if p.tokens[p.index].tokenType != LEFT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing left parenthesis")
			os.Exit(65)
		}
		p.index++
		expr, err := p.parseAssignment()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		if p.tokens[p.index].tokenType != RIGHT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing right parenthesis")
			os.Exit(65)
		}
		p.index++
		isBlock := false
		if p.tokens[p.index].tokenType == LEFT_BRACE {
			p.index++
			isBlock = true
		}

		statements := make([]Statement, 0)
		if isBlock {
			for p.index < len(p.tokens) &&
				p.tokens[p.index].tokenType != RIGHT_BRACE &&
				p.tokens[p.index].tokenType != EOF &&
				p.tokens[p.index].tokenType != ELSE {

				statement := p.parseStatement()
				if statement == nil {
					fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
					panic("error while parsing")
				}
				statements = append(statements, statement)
			}
			if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
				fmt.Fprintln(os.Stderr, "Missing right brace")
				os.Exit(65)
			}
			p.index++
		} else {
			statement := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}

		elseStatements := make([]Statement, 0)
		elseIfStatements := make([]IfStatement, 0)
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType == ELSE {
			p.index++
			elseBlock := false
			if p.index >= len(p.tokens) {
				fmt.Fprintln(os.Stderr, "Missing else block")
				os.Exit(65)
			}
			if p.tokens[p.index].tokenType == IF {
				// else if の場合
				p.index++

				if p.tokens[p.index].tokenType != LEFT_PAREN {
					fmt.Fprintln(os.Stderr, "Missing left parenthesis")
					os.Exit(65)
				}
				p.index++
				elseIfExpr, err := p.parseAssignment()
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(65)
				}
				if p.tokens[p.index].tokenType != RIGHT_PAREN {
					fmt.Fprintln(os.Stderr, "Missing right parenthesis")
					os.Exit(65)
				}
				p.index++
				isBlock := false
				if p.tokens[p.index].tokenType == LEFT_BRACE {
					p.index++
					isBlock = true
				}
				tmpStatements := make([]Statement, 0)
				if isBlock {
					for p.index < len(p.tokens) &&
						p.tokens[p.index].tokenType != RIGHT_BRACE &&
						p.tokens[p.index].tokenType != EOF &&
						p.tokens[p.index].tokenType != ELSE {
						statement := p.parseStatement()
						if statement == nil {
							fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
							panic("error while parsing")
						}
						tmpStatements = append(tmpStatements, statement)
					}
					if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
						fmt.Fprintln(os.Stderr, "Missing right brace")
						os.Exit(65)
					}
					p.index++
				} else {
					if p.index < len(p.tokens) &&
						p.tokens[p.index].tokenType != RIGHT_BRACE &&
						p.tokens[p.index].tokenType != EOF &&
						p.tokens[p.index].tokenType != ELSE {
						statement := p.parseStatement()
						if statement == nil {
							fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
							panic("error while parsing")
						}
						tmpStatements = append(tmpStatements, statement)
					}
				}
				elseIfStatements = append(elseIfStatements, IfStatement{
					expr:       elseIfExpr,
					statements: tmpStatements,
				})
			} else {
				//ただの else の場合
				if p.tokens[p.index].tokenType == LEFT_BRACE {
					p.index++
					elseBlock = true
					for p.index < len(p.tokens) &&
						p.tokens[p.index].tokenType != RIGHT_BRACE &&
						p.tokens[p.index].tokenType != EOF &&
						p.tokens[p.index].tokenType != ELSE {
						statement := p.parseStatement()
						if statement == nil {
							fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
							panic("error while parsing")
						}
						elseStatements = append(elseStatements, statement)
					}
				} else {
					if p.index < len(p.tokens) && p.tokens[p.index].tokenType != SEMICOLON && p.tokens[p.index].tokenType != EOF {
						statement := p.parseStatement()
						if statement == nil {
							fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
							panic("error while parsing")
						}
						elseStatements = append(elseStatements, statement)
					}
				}
				if elseBlock {
					if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
						fmt.Fprintln(os.Stderr, "Missing right brace")
						os.Exit(65)
					}
					p.index++
				}
			}
		}

		return &IfStatement{
			expr:             expr,
			statements:       statements,
			elseStatements:   elseStatements,
			elseIfStatements: elseIfStatements,
		}
	} else if p.tokens[p.index].tokenType == LEFT_BRACE {
		p.index++
		statements := make([]Statement, 0)
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType != RIGHT_BRACE && p.tokens[p.index].tokenType != EOF {
			statement := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}
		if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
			fmt.Fprintln(os.Stderr, "Missing right brace")
			os.Exit(65)
		}
		p.index++
		return &BlockStatement{
			statements: statements,
		}
	} else if p.tokens[p.index].tokenType == PRINT {
		p.index++
		if p.tokens[p.index].value == ";" {
			fmt.Fprintln(os.Stderr, "Missing expression after print")
			os.Exit(65)
		}
		expr, err := p.parseAssignment()
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
	} else if p.tokens[p.index].tokenType == VAR {
		p.index++
		varName := p.tokens[p.index].value
		p.index++
		if p.tokens[p.index].tokenType == SEMICOLON {
			p.index++

			return &VariableStatement{
				expr:    &NilNode{value: "nil", tokenType: NIL},
				varName: varName,
			}
		}
		if p.tokens[p.index].value != "=" {
			panic("= is missing")
		}
		p.index++
		varValue, err := p.parseAssignment()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		if p.tokens[p.index].value != ";" {
			panic("; is missing " + p.tokens[p.index].value)
		}
		p.index++

		return &VariableStatement{
			expr:    varValue,
			varName: varName,
		}
	} else if p.tokens[p.index].tokenType == WHILE {
		p.index++
		if p.tokens[p.index].tokenType != LEFT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing left parenthesis")
			os.Exit(65)
		}
		p.index++
		expr, err := p.parseAssignment()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		if p.tokens[p.index].tokenType != RIGHT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing right parenthesis")
			os.Exit(65)
		}
		p.index++
		isBlock := false
		if p.tokens[p.index].tokenType == LEFT_BRACE {
			p.index++
			isBlock = true
		}

		statements := make([]Statement, 0)
		if isBlock {
			for p.index < len(p.tokens) &&
				p.tokens[p.index].tokenType != RIGHT_BRACE &&
				p.tokens[p.index].tokenType != EOF {

				statement := p.parseStatement()
				if statement == nil {
					fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
					panic("error while parsing")
				}
				statements = append(statements, statement)
			}
			if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
				fmt.Fprintln(os.Stderr, "Missing right brace")
				os.Exit(65)
			}
			p.index++
		} else {
			statement := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}

		return &WhileStatement{
			expr:       expr,
			statements: statements,
		}
	} else if p.tokens[p.index].tokenType == FOR {
		p.index++;
		if p.tokens[p.index].tokenType != LEFT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing left parenthesis")
			os.Exit(65)
		}

		p.index++

		var firstStatement Statement
		// セミコロンでなければ、最初の文をパースする
		if p.tokens[p.index].tokenType != SEMICOLON {
			firstStatement = p.parseStatement()
		} else {
			// セミコロンの場合は、nil を代入する
			p.index++
			firstStatement = &ExpressionStatement{
				expr: &NilNode{value: "nil", tokenType: NIL},
			}
		}
		expression, err := p.parseAssignment()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		p.index++
		var endStatement Statement
		if p.tokens[p.index].tokenType != RIGHT_PAREN {
			expr, err := p.parseAssignment()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(65)
			}
			endStatement = &ExpressionStatement{expr: expr}
			if p.tokens[p.index].tokenType == SEMICOLON {
				p.index++
			}
			p.index++
		} else {
			// セミコロンの場合は、nil を代入する
			p.index++
			endStatement = &ExpressionStatement{
				expr: &NilNode{value: "nil", tokenType: NIL},
			}
		}

		isBlock := false

		if p.tokens[p.index].tokenType == LEFT_BRACE {
			p.index++
			isBlock = true
		}
		statements := make([]Statement, 0)
		if isBlock {
			for p.index < len(p.tokens) &&
				p.tokens[p.index].tokenType != RIGHT_BRACE &&
				p.tokens[p.index].tokenType != EOF {

				statement := p.parseStatement()
				if statement == nil {
					fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
					panic("error while parsing")
				}
				statements = append(statements, statement)
			}
			if p.index >= len(p.tokens) || p.tokens[p.index].tokenType != RIGHT_BRACE {
				fmt.Fprintln(os.Stderr, "Missing right brace")
				os.Exit(65)
			}
			p.index++
		} else {
			statement := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}

		return &ForStatement{
			firstStatement:  firstStatement,
			expression:      expression,
			endStatement:    endStatement,
			statements:      statements,
		}
	}

	// ただの式。特に何かをしているわけではない。
	expression, err := p.parseAssignment()
	if p.tokens[p.index].tokenType == SEMICOLON {
		p.index++
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(65)
		}
		return &ExpressionStatement{
			expr: expression,
		}
	}

	// fmt.Printf("Unknown token: %s\n", p.tokens[p.index].tokenType)
	panic("unknown statement")
}

func (p *Parser) parseAssignment() (Node, error) {
	if p.index >= len(p.tokens) {
		panic("Index out of range")
	}
	token := p.tokens[p.index]

	if token.tokenType == IDENTIFIER {
		if p.tokens[p.index+1].value == SEMICOLON {
			p.index++
			return &IdentifierNode{
				value:     token.value,
				tokenType: token.tokenType,
			}, nil
		}
		if p.tokens[p.index+1].tokenType == EQUAL {
			p.index++
			p.index++
			value, err := p.parseAssignment()
			if err != nil {
				return nil, err
			}
			return &AssignmentNode{
				varName:   token.value,
				value:     value,
				valueType: ASSIGNMENT,
			}, nil
		}
	}

	node, err := p.parseExpression()
	// Identifier でない場合は expression をそのまま返す

	if p.tokens[p.index].tokenType == OR {
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType == OR {
			or_token := p.tokens[p.index]
			p.index++
			rightNode, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if rightNode == nil {
				fmt.Fprintln(os.Stderr, "Error parsing right node")
				os.Exit(65)
			}
			node = &Binary{
				left:      node,
				operator:  or_token,
				right:     rightNode,
				tokenType: OR,
			}
		}
		return node, nil
	} else if p.tokens[p.index].tokenType == AND {
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType == AND {
			and_token := p.tokens[p.index]
			p.index++
			rightNode, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			node = &Binary{
				left:      node,
				operator:  and_token,
				right:     rightNode,
				tokenType: AND,
			}
		}
		return node, nil
	}

	return node, err
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
		expression, _ = p.parseAssignment()
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

	// fmt.Printf("Unknown token: %s\n", token.tokenType)
	return &NilNode{}, fmt.Errorf("unknown expression")
}

func (u *Unary) getValue(env *Env) EvaluateNode {
	if u.operator.tokenType == MINUS {
		if u.right.getValue(env).valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operand must be a number.")
			os.Exit(70)
		}
		num, _ := strconv.Atoi(u.right.getValue(env).value)
		return EvaluateNode{
			value:     strconv.FormatFloat(float64(num*-1), 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if u.operator.tokenType == BANG {
		if u.right.getValue(env).value == "false" || u.right.getValue(env).value == "" || u.right.getValue(env).value == "nil" {
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

func (g *Group) getValue(env *Env) EvaluateNode {
	if len(g.nodes) == 1 {
		return EvaluateNode{
			value:     g.nodes[0].getValue(env).value,
			valueType: g.nodes[0].getValue(env).valueType,
		}
	} else {
		values := ""
		for i, n := range g.nodes {
			if i != 0 {
				values += " "
			}
			values += n.getValue(env).value
		}
		return EvaluateNode{
			value:     values,
			valueType: STRING,
		}
	}
}

func (a *AssignmentNode) getValue(env *Env) EvaluateNode {
	variables := (*env.variables)
	if _, ok := variables[a.varName]; !ok {
		fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n", a.varName)
		os.Exit(70)
	}
	// 変数の値をセットする
	value := a.value.getValue(env).value
	valueType := a.value.getValue(env).valueType
	newValue := EvaluateNode{
		value:     value,
		valueType: valueType,
	}
	variables[a.varName] = newValue

	if _, ok := (*env.parentVariables)[a.varName]; ok {
		(*env.parentVariables)[a.varName] = newValue
	}

	// 変数の値を返す
	return EvaluateNode{
		value:     value,
		valueType: valueType,
	}
}

func (i *IdentifierNode) getValue(env *Env) EvaluateNode {
	variables := *env.variables

	if _, ok := variables[i.value]; !ok {
		fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n", i.value)
		os.Exit(70)
	}

	return EvaluateNode{
		value:     variables[i.value].value,
		valueType: variables[i.value].valueType,
	}
}

func (s *StringNode) getValue(env *Env) EvaluateNode {
	return EvaluateNode{
		value:     s.value,
		valueType: STRING,
	}
}
func (n *NumberNode) getValue(env *Env) EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NUMBER,
	}
}
func (b *BooleanNode) getValue(env *Env) EvaluateNode {
	return EvaluateNode{
		value:     b.value,
		valueType: BOOLEAN,
	}
}
func (n *NilNode) getValue(env *Env) EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NIL,
	}
}

func (b *Binary) getValue(env *Env) EvaluateNode {
	if b.operator.tokenType == PLUS {
		if b.left.getValue(env).valueType != b.right.getValue(env).valueType {
			fmt.Fprintf(os.Stderr, "Operands must be same types.")
			os.Exit(70)
		}
		if b.left.getValue(env).valueType == STRING {
			return EvaluateNode{
				value:     b.left.getValue(env).value + b.right.getValue(env).value,
				valueType: STRING,
			}
		} else if b.left.getValue(env).valueType == NUMBER {
			left, _ := strconv.ParseFloat(b.left.getValue(env).value, 10)
			right, _ := strconv.ParseFloat(b.right.getValue(env).value, 10)
			return EvaluateNode{
				value:     strconv.FormatFloat(left+right, 'f', -1, 64),
				valueType: NUMBER,
			}
		}
	}

	left, _ := strconv.ParseFloat(b.left.getValue(env).value, 10)
	right, _ := strconv.ParseFloat(b.right.getValue(env).value, 10)
	if b.operator.tokenType == SLASH {
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left/right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == STAR {
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left*right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == MINUS {
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operands must be numbers.")
			os.Exit(70)
		}
		return EvaluateNode{
			value:     strconv.FormatFloat(left-right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == GREATER {
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
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
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
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
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
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
		if b.left.getValue(env).valueType != NUMBER || b.right.getValue(env).valueType != NUMBER {
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
		if b.left.getValue(env).value == b.right.getValue(env).value && b.left.getValue(env).valueType == b.right.getValue(env).valueType {
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
		if b.left.getValue(env).value != b.right.getValue(env).value || b.left.getValue(env).valueType != b.right.getValue(env).valueType {
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

	if b.operator.tokenType == OR {
		leftValue := b.left.getValue(env)
		if isTrueString(leftValue.value) {
			return leftValue
		}

		rightValue := b.right.getValue(env)
		if isTrueString(rightValue.value) {
			return rightValue
		}

		return EvaluateNode{
			value:     "false",
			valueType: BOOLEAN,
		}
	} else if b.operator.tokenType == AND {
		if isTrueString(b.left.getValue(env).value) && isTrueString(b.right.getValue(env).value) {
			return EvaluateNode{
				value:     b.right.getValue(env).value,
				valueType: b.right.getValue(env).valueType,
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

func (e *EvaluateNode) getValue(env *Env) EvaluateNode {
	return EvaluateNode{
		value:     e.getValue(env).value,
		valueType: e.getValue(env).valueType,
	}
}

func isTrueString(value string) bool {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	// から文字は true とみなす
	// if value == "" {
	// 	return false
	// }
	if value == "nil" {
		return false
	}
	return true
}
