package run

import (
	"fmt"
)

type Statement interface {
	Execute(env *Env) *ReturnError
}

type PrintStatement struct {
	Statement
	expr Node
}

type BlockStatement struct {
	Statement
	statements []Statement
}

type FunStatement struct {
	Statement
	name       string
	parameters []string
	statements []Statement
	closure    *Env
}

type ExpressionStatement struct {
	Statement
	expr Node
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
type WhileStatement struct {
	Statement
	expr       Node
	statements []Statement
}

type ForStatement struct {
	Statement
	firstStatement Statement
	expression     Node
	endStatement   Statement
	// for の中の文
	statements []Statement
}

type ReturnStatement struct {
	Statement
	expr Node
}

type ReturnError struct {
	value     string
	valueType string
	function  *Function
}

func (e *ExpressionStatement) Execute(env *Env) *ReturnError {
	e.expr.getValue(env)
	return nil
}

func (b *BlockStatement) Execute(env *Env) *ReturnError {
	newEnv := env.NewChildEnv()
	for _, statement := range b.statements {
		if err := statement.Execute(newEnv); err != nil {
			return err
		}
	}

	// ブロック終了時に親環境を更新
	return nil
}

func (p *PrintStatement) Execute(env *Env) *ReturnError {
	value := p.expr.getValue(env).value
	fmt.Println(value)

	return nil
}

func (f *ForStatement) Execute(parentEnv *Env) *ReturnError {
	newEnv := parentEnv.NewChildEnv()

	if err := f.firstStatement.Execute(newEnv); err != nil {
		return err
	}

	for isTrueString(f.expression.getValue(newEnv).value) {
		grandChildEnv := newEnv.NewChildEnv()
		for _, statement := range f.statements {
			if err := statement.Execute(grandChildEnv); err != nil {
				return err
			}
		}
		if err := f.endStatement.Execute(newEnv); err != nil {
			return err
		}
	}
	return nil
}

func (v *VariableStatement) Execute(env *Env) *ReturnError {
	value := v.expr.getValue(env)
	// 新しい変数を現在の環境に定義
	env.Define(v.varName, value)
	return nil
}

func (i *IfStatement) Execute(parentEnv *Env) *ReturnError {
	value := i.expr.getValue(parentEnv)
	newEnv := parentEnv.NewChildEnv()
	statements := []Statement{}
	if isTrueString(value.value) {
		statements = i.statements
	} else if len(i.elseIfStatements) > 0 {
		for _, elseIfStatement := range i.elseIfStatements {
			// 何も条件に引っ掛からなかった場合は else を実行する
			statements = i.elseStatements

			if isTrueString(elseIfStatement.expr.getValue(parentEnv).value) {
				statements = elseIfStatement.statements
				break
			}
		}
	} else {
		statements = i.elseStatements
	}

	for _, statement := range statements {
		if err := statement.Execute(newEnv); err != nil {
			return err
		}
	}
	return nil
}

func (w *WhileStatement) Execute(parentEnv *Env) *ReturnError {
	newEnv := parentEnv.NewChildEnv()
	for isTrueString(w.expr.getValue(newEnv).value) {
		if len(w.statements) > 0 {
			for _, statement := range w.statements {
				if err := statement.Execute(newEnv); err != nil {
					return err
				}
			}
		}
		// ループ内での変更を親環境に反映
	}
	return nil
}


func (f *FunStatement) Execute(env *Env) *ReturnError {
	// 関数を定義
	fn := Function{
		name:       f.name,
		parameters: f.parameters,
		statements: f.statements,
		closure:    env, // 現在の環境をクロージャとして保持
	}

	// 関数を変数として定義
	env.Define(f.name, EvaluateNode{
		value:     "<fn " + f.name + ">",
		valueType: "function",
		function:  &fn,
	})

	return nil
}

func (r *ReturnStatement) Execute(env *Env) *ReturnError {
	node := r.expr.getValue(env)

	if node.valueType == "function" {
		return &ReturnError{
			value:     node.value,
			valueType: "function",
			function:  node.function,
		}
	}

	return &ReturnError{
		value:     node.value,
		valueType: node.valueType,
	}
}


func (r *ReturnError) Error() string {
	return r.valueType + " " + r.value
}
