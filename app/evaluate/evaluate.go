package evaluate

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func Evaluate() {
	ast := Parse()
	ast.GetValue()
}

// Print methods for AST and nodes
func (a *AST) GetValue() {
	for _, n := range a.nodes {
		io.WriteString(os.Stdout, n.getValue().value)
		io.WriteString(os.Stdout, "\n")
	}
}

func (u *Unary) getValue() EvaluateNode {
	if u.operator.tokenType == MINUS {
		if u.right.getValue().valueType != NUMBER {
			fmt.Fprintf(os.Stderr, "Operand must be a number.")
			os.Exit(70)
		}
		num, _ := strconv.Atoi(u.right.getValue().value)
		return EvaluateNode{
			value:     strconv.FormatFloat(float64(num*-1), 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if u.operator.tokenType == BANG {
		if u.right.getValue().value == "false" || u.right.getValue().value == "" || u.right.getValue().value == "nil" {
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

func (g *Group) getValue() EvaluateNode {
	if len(g.nodes) == 1 {
		return EvaluateNode{
			value:     g.nodes[0].getValue().value,
			valueType: g.nodes[0].getValue().valueType,
		}
	} else {
		values := ""
		for i, n := range g.nodes {
			if i != 0 {
				values += " "
			}
			values += n.getValue().value
		}
		return EvaluateNode{
			value:     values,
			valueType: STRING,
		}
	}
}

func (s *StringNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     s.value,
		valueType: STRING,
	}
}
func (n *NumberNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NUMBER,
	}
}
func (b *BooleanNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     b.value,
		valueType: BOOLEAN,
	}
}
func (n *NilNode) getValue() EvaluateNode {
	return EvaluateNode{
		value:     n.value,
		valueType: NIL,
	}
}

func (b *Binary) getValue() EvaluateNode {
	if b.operator.tokenType == SLASH {
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
		return EvaluateNode{
			value:     strconv.FormatFloat(left/right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == STAR {
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
		return EvaluateNode{
			value:     strconv.FormatFloat(left*right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == MINUS {
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
		return EvaluateNode{
			value:     strconv.FormatFloat(left-right, 'f', -1, 64),
			valueType: NUMBER,
		}
	} else if b.operator.tokenType == PLUS {
		if b.left.getValue().valueType == STRING {
			return EvaluateNode{
				value:     b.left.getValue().value + b.right.getValue().value,
				valueType: STRING,
			}
		} else if b.left.getValue().valueType == NUMBER {
			left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
			right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
			return EvaluateNode{
				value:     strconv.FormatFloat(left+right, 'f', -1, 64),
				valueType: NUMBER,
			}
		}
	} else if b.operator.tokenType == GREATER {
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
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
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
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
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
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
		left, _ := strconv.ParseFloat(b.left.getValue().value, 10)
		right, _ := strconv.ParseFloat(b.right.getValue().value, 10)
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
		if b.left.getValue().value == b.right.getValue().value && b.left.getValue().valueType == b.right.getValue().valueType {
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
		if b.left.getValue().value != b.right.getValue().value || b.left.getValue().valueType != b.right.getValue().valueType {
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

	panic("Unknown operator: " + b.operator.tokenType)
}
