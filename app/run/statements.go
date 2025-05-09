package run

import (
	"fmt"
)

type Statement interface {
	Execute(env *Env) error
}

type PrintStatement struct {
	Statement
	expr Node
}

type BlockStatement struct {
	Statement
	statements []Statement
}

func (b *BlockStatement) Execute(env *Env) error {
	newEnv := env.NewChildEnv()
	for _, statement := range b.statements {
		if err := statement.Execute(newEnv); err != nil {
			return err
		}
	}
	return nil
}

func (p *PrintStatement) Execute(env *Env) error {
	value := p.expr.getValue(env).value
	fmt.Println(value)

	return nil
}

type ExpressionStatement struct {
	Statement
	expr Node
}

func (e *ExpressionStatement) Execute(env *Env) error {
	e.expr.getValue(env)
	return nil
}

type VariableStatement struct {
	Statement
	expr    Node
	varName string
}

func (v *VariableStatement) Execute(env *Env) error {
	value := v.expr.getValue(env)
	// 変数の値をセットする
	env.variables[v.varName] = value
	return nil
}
