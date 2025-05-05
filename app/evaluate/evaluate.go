package evaluate

import (
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
		io.WriteString(os.Stdout, n.getValue())
		io.WriteString(os.Stdout, "\n")
	}
}

func (u *Unary) getValue() string {
	if u.operator.tokenType == "MINUS" {
		return "-" + u.right.getValue()
	} else if u.operator.tokenType == "BANG" {
		if u.right.getValue() == "false" || u.right.getValue() == "" || u.right.getValue() == "nil" {
			return "true"
		} else {
			return "false"
		}
	}

	panic("Unknown operator: " + u.operator.tokenType)
}

func (g *Group) getValue() string {
	values := ""
	for i, n := range g.nodes {
		if i != 0 {
			values += " "
		}
		values += n.getValue()
	}
	return values
}
func (s *StringNode) getValue() string {
	return s.value
}
func (n *NumberNode) getValue() string {
	return n.value
}
func (b *BooleanNode) getValue() string {
	return b.value
}
func (n *NilNode) getValue() string {
	return n.value
}

// binary については一旦考えない
func (b *Binary) getValue() string {
	left, _ := strconv.Atoi(b.left.getValue())
	right, _ := strconv.Atoi(b.right.getValue())
	if b.operator.tokenType == "SLASH" {
		return strconv.Itoa(left / right)
	} else if b.operator.tokenType == "STAR"{
		return strconv.Itoa(left * right)
	}

	panic("Unknown operator: " + b.operator.tokenType)
}
