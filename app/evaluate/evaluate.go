package evaluate

import (
	"io"
	"os"
)

func Evaluate() {
	ast := Parse()
	ast.Print()
}

// Print methods for AST and nodes
func (a *AST) Print() {
	for _, n := range a.nodes {
		n.Print()
		io.WriteString(os.Stdout, "\n")
	}
}

func (u *Unary) Print() {
	io.WriteString(os.Stdout, "(")
	io.WriteString(os.Stdout, u.operator.value)
	io.WriteString(os.Stdout, " ")
	u.right.Print()
	io.WriteString(os.Stdout, ")")
}

func (g *Group) Print() {
	for _, n := range g.nodes {
		n.Print()
	}
}
func (s *StringNode) Print() {
	io.WriteString(os.Stdout, s.value)
}
func (n *NumberNode) Print() {
	io.WriteString(os.Stdout, n.value)
}
func (n *BooleanNode) Print() {
	io.WriteString(os.Stdout, n.value)
}
func (n *NilNode) Print() {
	io.WriteString(os.Stdout, n.value)
}

func (b *Binary) Print() {
	io.WriteString(os.Stdout, "(")
	io.WriteString(os.Stdout, b.operator.value)
	io.WriteString(os.Stdout, " ")
	b.left.Print()
	io.WriteString(os.Stdout, " ")
	b.right.Print()
	io.WriteString(os.Stdout, ")")
}
