package run

type Token struct {
	tokenType string
	value     string
}

type Parser struct {
	tokens []Token
	index  int
}

type Node interface {
	getValue(env *Env) EvaluateNode
	getType() string
}

type EvaluateNode struct {
	value     string
	valueType string
	function  *Function // 関数の場合
}

type AssignmentNode struct {
	Node
	varName   string
	value     Node
	valueType string
}

type StringNode struct {
	Node
	value     string
	tokenType string
}

type NumberNode struct {
	Node
	value     string
	tokenType string
}

type BooleanNode struct {
	Node
	value     string
	tokenType string
}

type NilNode struct {
	Node
	value     string
	tokenType string
}

type Group struct {
	Node
	nodes     []Node
	tokenType string
}

type Unary struct {
	Node
	operator  Token
	right     Node
	tokenType string
}

type Binary struct {
	Node
	left      Node
	operator  Token
	right     Node
	tokenType string
}

type IdentifierNode struct {
	Node
	value     string
	tokenType string
}

type FuncNode struct {
	Node
	callee    Node
	arguments []Node
	tokenType string
}

func (s *StringNode) getType() string {
	return s.tokenType
}

func (n *NumberNode) getType() string {
	return n.tokenType
}

func (b *BooleanNode) getType() string {
	return b.tokenType
}

func (n *NilNode) getType() string {
	return n.tokenType
}

func (g *Group) getType() string {
	return g.tokenType
}

func (u *Unary) getType() string {
	return u.tokenType
}

func (b *Binary) getType() string {
	return b.tokenType
}

func (i *IdentifierNode) getType() string {
	return i.tokenType
}

func (f *FuncNode) getType() string {
	return f.tokenType
}

func (a *AssignmentNode) getType() string {
	return a.valueType
}