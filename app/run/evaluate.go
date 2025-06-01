package run

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func (u *Unary) getValue(env *Env) EvaluateNode {
	if u.operator.tokenType == MINUS {
		if u.right.getValue(env).valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operand must be a number.")
			os.Exit(70)
		}
		num, _ := strconv.ParseFloat(u.right.getValue(env).value, 64)
		return EvaluateNode{
			value:     strconv.FormatFloat(-num, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if u.operator.tokenType == BANG {
		value := u.right.getValue(env).value
		if value == "false" || value == "" || value == "nil" {
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
	// 値を評価
	result := a.value.getValue(env)
	
	// 変数に値をセット
	if !env.Set(a.varName, result) {
		fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n", a.varName)
		os.Exit(70)
	}
	
	return result
}

func (i *IdentifierNode) getValue(env *Env) EvaluateNode {
	// 特殊な組み込み関数の場合
	if i.value == "clock" {
		return EvaluateNode{
			value:     "clock",
			valueType: STRING,
		}
	}

	// 変数を探す
	if val, ok := env.Get(i.value); ok {
		return val
	}

	fmt.Fprintf(os.Stderr, "Undefined variable '%s'.\n", i.value)
	os.Exit(70)
	return EvaluateNode{}
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
		value:     e.value,
		valueType: e.valueType,
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

func (f *FuncNode) getValue(env *Env) EvaluateNode {
	// calleeを評価
	calleeValue := f.callee.getValue(env)
	
	// clock関数の特別処理
	if calleeValue.value == "clock" {
		return EvaluateNode{
			value:     strconv.FormatInt(time.Now().Unix(), 10),
			valueType: NUMBER,
		}
	}

	// 関数を取得
	var funcDef *Function
	if calleeValue.valueType == "function" && calleeValue.function != nil {
		// calleeが直接関数値を返した場合
		funcDef = calleeValue.function
	} else {
		// calleeが関数名を返した場合（後方互換性のため）
		if identNode, ok := f.callee.(*IdentifierNode); ok {
			if val, ok := env.Get(identNode.value); ok && val.valueType == "function" {
				funcDef = val.function
			}
		}
	}

	// 関数が見つからなかったらエラー
	if funcDef == nil {
		fmt.Fprintf(os.Stderr, "Undefined function '%s'.\n", calleeValue.value)
		os.Exit(70)
	}

	if len(f.arguments) != len(funcDef.parameters) {
		fmt.Fprintf(os.Stderr, "Function '%s' expects %d arguments, but got %d.\n", funcDef.name, len(funcDef.parameters), len(f.arguments))
		os.Exit(70)
	}
	// 関数のクロージャ環境から新しい環境を作成
	newEnv := funcDef.closure.NewChildEnv()

	// 引数を新しい環境にバインド
	for index, arg := range f.arguments {
		argument := arg.getValue(env)
		newEnv.Define(funcDef.parameters[index], argument)
	}

	for _, statement := range funcDef.statements {
		// 実際にはエラーではないが、エラーとして扱う
		// 実際には return で返ってくるものが入っている
		err := statement.Execute(newEnv)
		if err != nil {
			return EvaluateNode{
				value:     err.value,
				valueType: err.valueType,
				function:  err.function,
			}
		}
	}

	return EvaluateNode{
		value:     "nil",
		valueType: NIL,
	}
}