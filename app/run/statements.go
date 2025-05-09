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

type BlockStatement struct {
	Statement
	statements []Statement
}

func (b *BlockStatement) Execute() error {
	for _, statement := range b.statements {
		if err := statement.Execute(); err != nil {
			return err
		}
	}
	return nil
}

func (p *PrintStatement) Execute() error {
	value := p.expr.getValue().value
	fmt.Println(value)

	return nil
}

type ExpressionStatement struct {
	Statement
	expr Node
}

func (e *ExpressionStatement) Execute() error {
	e.expr.getValue()
	return nil
}

type VariableStatement struct {
	Statement
	expr    Node
	varName string
}

func (v *VariableStatement) Execute() error {
	globalEnv := getGlobalEnv()
	value := v.expr.getValue()
	// 変数の値をセットする
	globalEnv.variables[v.varName] = value
	return nil
}
