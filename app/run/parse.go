package run

import (
	"fmt"
	"os"
)



// parse して構文木を作成する
func (p *Parser) parseStatements() []Statement {
	statements := make([]Statement, 0)
	for p.index < len(p.tokens) && p.tokens[p.index].tokenType != EOF {
		statement, err := p.parseStatement()
		if err != nil {
			fmt.Fprintln(os.Stderr, err) // Changed from os.Err(err) to fmt.Fprintln
			os.Exit(65)
		}
		if statement == nil {
			fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
			panic("error while parsing")
		}
		statements = append(statements, statement)
	}
	return statements
}

func (p *Parser) parseStatement() (Statement, error) {
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

				statement, _ := p.parseStatement()
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
			statement, _ := p.parseStatement()
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
						statement, _ := p.parseStatement()
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
						statement, _ := p.parseStatement()
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
						statement, _ := p.parseStatement()
						if statement == nil {
							fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
							panic("error while parsing")
						}
						elseStatements = append(elseStatements, statement)
					}
				} else {
					if p.index < len(p.tokens) && p.tokens[p.index].tokenType != SEMICOLON && p.tokens[p.index].tokenType != EOF {
						statement, _ := p.parseStatement()
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
		}, nil
	} else if p.tokens[p.index].tokenType == LEFT_BRACE {
		p.index++
		statements := make([]Statement, 0)
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType != RIGHT_BRACE && p.tokens[p.index].tokenType != EOF {
			statement, _ := p.parseStatement()
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
		}, nil
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
			}, nil
		}

		return nil, fmt.Errorf("syntax error")
	} else if p.tokens[p.index].tokenType == VAR {
		p.index++
		varName := p.tokens[p.index].value
		p.index++
		if p.tokens[p.index].tokenType == SEMICOLON {
			p.index++

			return &VariableStatement{
				expr:    &NilNode{value: "nil", tokenType: NIL},
				varName: varName,
			}, nil
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
		}, nil
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

				statement, _ := p.parseStatement()
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
			statement, _ := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}

		return &WhileStatement{
			expr:       expr,
			statements: statements,
		}, nil
	} else if p.tokens[p.index].tokenType == FOR {
		p.index++
		if p.tokens[p.index].tokenType != LEFT_PAREN {
			fmt.Fprintln(os.Stderr, "Missing left parenthesis")
			os.Exit(65)
		}

		p.index++

		var firstStatement Statement
		// セミコロンでなければ、最初の文をパースする
		if p.tokens[p.index].tokenType != SEMICOLON {
			firstStatement, _ = p.parseStatement()
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

				statement, _ := p.parseStatement()
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
			statement, _ := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}

		return &ForStatement{
			firstStatement: firstStatement,
			expression:     expression,
			endStatement:   endStatement,
			statements:     statements,
		}, nil
	} else if p.tokens[p.index].tokenType == FUN {
		// fun の部分の index を ++ する
		p.index++

		funName := p.tokens[p.index].value

		p.index++

		if p.tokens[p.index].tokenType != LEFT_PAREN {
			// 一旦 syntax error にしておく
			return nil, fmt.Errorf("syntax error")
		}

		p.index++
		var parameters []string
		for p.tokens[p.index].tokenType != RIGHT_PAREN {
			parameter := p.tokens[p.index].value
			parameters = append(parameters, parameter)
			p.index++
			if p.tokens[p.index].tokenType != COMMA {
				break
			}
			p.index++
		}
		p.index++

		if p.tokens[p.index].tokenType != LEFT_BRACE {
			return nil, fmt.Errorf("syntax error")
		}

		p.index++
		statements := make([]Statement, 0)
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType != RIGHT_BRACE {
			statement, _ := p.parseStatement()
			if statement == nil {
				fmt.Fprintf(os.Stderr, "Error parsing statement at index %d\n", p.index)
				panic("error while parsing")
			}
			statements = append(statements, statement)
		}
		// to do
		p.index++
		return &FunStatement{
			name:       funName,
			parameters: parameters,
			statements: statements,
		}, nil
	} else if p.tokens[p.index].tokenType == RETURN {
		p.index++
		var expr Node
		var err error
		if p.tokens[p.index].tokenType == SEMICOLON {
			p.index++
			expr = &NilNode{value: "nil", tokenType: NIL}
		} else {
			expr, err = p.parseAssignment()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(65)
			}
			if p.tokens[p.index].tokenType == SEMICOLON {
				p.index++
			}
		}
		return &ReturnStatement{
			expr: expr,
		}, nil
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
		}, nil
	}

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

	return p.parseCall()
}

func (p *Parser) parseCall() (Node, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.index < len(p.tokens) && p.tokens[p.index].tokenType == LEFT_PAREN {
		p.index++
		args := make([]Node, 0)
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType != RIGHT_PAREN {
			arg, err := p.parseAssignment()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.index < len(p.tokens) && p.tokens[p.index].tokenType == COMMA {
				p.index++
			}
		}
		if p.index < len(p.tokens) && p.tokens[p.index].tokenType == RIGHT_PAREN {
			p.index++
		} else {
			return nil, fmt.Errorf("missing right parenthesis")
		}

		chain := false
		for p.index < len(p.tokens) && p.tokens[p.index].tokenType == LEFT_PAREN {
			chain = true
			expr = &FuncNode{
				callee:    expr,
				arguments: args,
				tokenType: FUN,
			}
			p.index++
			args = make([]Node, 0)
			for p.index < len(p.tokens) && p.tokens[p.index].tokenType != RIGHT_PAREN {
				arg, err := p.parseAssignment()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)

				if p.index < len(p.tokens) && p.tokens[p.index].tokenType == COMMA {
					p.index++
				}
			}
			if p.index < len(p.tokens) && p.tokens[p.index].tokenType == RIGHT_PAREN {
				p.index++
			} else {
				return nil, fmt.Errorf("missing right parenthesis")
			}

			expr = &FuncNode{
				callee:    expr,
				arguments: args,
				tokenType: FUN,
			}
		}

		// 関数チェーンの場合はそのまま返す
		if chain {
			return expr, nil
		}

		return &FuncNode{
			callee:    expr,
			arguments: args,
			tokenType: FUN,
		}, nil
	}

	return expr, nil
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

	return &NilNode{}, fmt.Errorf("unknown expression")
}


