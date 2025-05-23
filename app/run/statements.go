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
		setParentEnv(newEnv)
		return err
	}

	for isTrueString(f.expression.getValue(newEnv).value) {
		grandChildEnv := newEnv.NewChildEnv()
		for _, statement := range f.statements {
			if err := statement.Execute(grandChildEnv); err != nil {
				// setParentEnv(grandChildEnv)
				return err
			}
			// setParentEnv(grandChildEnv)
		}
		if err := f.endStatement.Execute(newEnv); err != nil {
			setParentEnv(newEnv)
			return err
		}
	}
	setParentEnv(newEnv)
	return nil
}

func (v *VariableStatement) Execute(env *Env) *ReturnError {
	value := v.expr.getValue(env)
	// 変数の値をセットする
	(env.variables)[v.varName] = value
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
		setParentEnv(newEnv)
	}
	return nil
}

func (w *WhileStatement) Execute(parentEnv *Env) *ReturnError {
	newEnv := parentEnv.NewChildEnv()
	for isTrueString(w.expr.getValue(newEnv).value) {
		for _, statement := range w.statements {
			if err := statement.Execute(newEnv); err != nil {
				setParentEnv(newEnv)
				return err
			}
			setParentEnv(newEnv)
		}
		setParentEnv(newEnv)
	}
	setParentEnv(newEnv)
	return nil
}

var functionIdCounter = 0

func genFunctionId(baseStr string) string {
	functionIdCounter++
	return fmt.Sprintf("%s-%d", baseStr, functionIdCounter)
}

func getFunctionId(baseStr string) string {
	return fmt.Sprintf("%s-%d", baseStr, functionIdCounter)
}

func (f *FunStatement) Execute(env *Env) *ReturnError {
	// 関数の定義をセットする
	(env.functions)["<fn "+f.name+">"] = Function{
		name:       f.name,
		parameters: f.parameters,
		statements: f.statements,
		closure:    env.newClosureEnv(),
	}
	(env.variables)[f.name] = EvaluateNode{
		value:     "<fn " + f.name + ">",
		valueType: "function",
	}

	(funcGlobalEnv.functions)[genFunctionId("<fn "+f.name+">")] = Function{
		name:       f.name,
		parameters: f.parameters,
		statements: f.statements,
		closure:    env.newClosureEnv(),
	}

	return nil
}

func (r *ReturnStatement) Execute(env *Env) *ReturnError {
	node := r.expr.getValue(env)

	if node.valueType == "function" {
		return &ReturnError{
			value:     getFunctionId(node.value),
			valueType: "function",
		}
	}

	return &ReturnError{
		value:     node.value,
		valueType: node.valueType,
	}
}

func setParentEnv(childEnv *Env) {
	// for k, v := range childEnv.parentVariables {
	// 	parentEnv := childEnv.parentEnv
	// 	for _, ok := parentEnv.variables[k]; ok; {
	// 		if parentEnv.variables[k] != v {
	// 			parentEnv.variables[k] = v
	// 		} else if _, ok := parentEnv.parentVariables[k]; ok && parentEnv.parentVariables[k] != v {
	// 			parentEnv.parentVariables[k] = v
	// 		} else {
	// 			break;
	// 		}
	// 		parentEnv = parentEnv.parentEnv
	// 	}
	// }
	for k, v := range childEnv.variables {
		// この変数が親から継承したものかどうかを確認
		if originalValue, exists := childEnv.parentVariables[k]; exists {
				// 値が変更された場合のみ更新する
				if originalValue != v {
						childEnv.parentEnv.variables[k] = v
				}
		}
		// 子環境で新たに定義された変数は親環境に反映しない
}
}

func (r *ReturnError) Error() string {
	return r.valueType + " " + r.value
}
