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

func (p *PrintStatement) Execute() error {
	value := p.expr.getValue().value
	fmt.Println(value)

	return nil
}
