package run

import (
	"fmt"
)

type Statement interface {
	Execute() error
}

type PrintStatement struct {
	Statement
	expr Node
}

type ExpressionStatement struct {
	Statement
	expr Node
}

func (p *PrintStatement) Execute() error {
	value := p.expr.getValue().value
	fmt.Println(value)

	return nil
}

func (e *ExpressionStatement) Execute() error {
	e.expr.getValue()
	return nil
}
