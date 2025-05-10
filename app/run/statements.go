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
	setParentEnv(env, newEnv)
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

// var xxx = yyy; の時に生成されるやつ
type VariableStatement struct {
	Statement
	expr    Node
	varName string
}

// if xxx { } の時に生成されるやつ
type IfStatement struct {
	Statement
	expr             Node
	statements       []Statement
	elseStatements   []Statement
	elseIfStatements []IfStatement
}

func (v *VariableStatement) Execute(env *Env) error {
	value := v.expr.getValue(env)
	// 変数の値をセットする
	(*env.variables)[v.varName] = value
	return nil
}

func (i *IfStatement) Execute(parentEnv *Env) error {
	value := i.expr.getValue(parentEnv)
	newEnv := parentEnv.NewChildEnv()
	statements := []Statement{}
	if value.value == "true" {
		statements = i.statements
	} else if len(i.elseIfStatements) > 0 {
		for _, elseIfStatement := range i.elseIfStatements {
			// 何も条件に引っ掛からなかった場合は else を実行する
			statements = i.elseStatements

			if elseIfStatement.expr.getValue(parentEnv).value == "true" {
				statements = elseIfStatement.statements
				break
			}
		}
	} else {
		statements = i.elseStatements
	}

	for _, statement := range statements {
		if err := statement.Execute(newEnv); err != nil {
			setParentEnv(parentEnv, newEnv)
			panic(err)
		}
		setParentEnv(parentEnv, newEnv)
	}
	return nil
}

func setParentEnv(parentEnv *Env, childEnv *Env) {
	for k, v := range *childEnv.parentVariables {
		if _, ok := (*parentEnv.variables)[k]; ok {
			if (*parentEnv.variables)[k] != v {
				(*parentEnv.variables)[k] = v
				(*parentEnv.parentVariables)[k] = v
			}
		}
	}
}
