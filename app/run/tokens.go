package run

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	LEFT_BRACE    = "LEFT_BRACE"
	RIGHT_BRACE   = "RIGHT_BRACE"
	COMMA         = "COMMA"
	DOT           = "DOT"
	MINUS         = "MINUS"
	PLUS          = "PLUS"
	SEMICOLON     = "SEMICOLON"
	STAR          = "STAR"
	EQUAL         = "EQUAL"
	EQUAL_EQUAL   = "EQUAL_EQUAL"
	BANG_EQUAL    = "BANG_EQUAL"
	LESS          = "LESS"
	LESS_EQUAL    = "LESS_EQUAL"
	GREATER       = "GREATER"
	GREATER_EQUAL = "GREATER_EQUAL"
	SLASH         = "SLASH"
	BANG          = "BANG"
	STRING        = "STRING"
	NUMBER        = "NUMBER"
	NIL           = "NIL"
	TRUE          = "TRUE"
	FALSE         = "FALSE"
	BOOLEAN       = "BOOLEAN"
	PRINT         = "PRINT"
	EOF           = "EOF"
	VAR           = "VAR"
	IDENTIFIER    = "IDENTIFIER"
	ASSIGNMENT    = "ASSIGNMENT"
	IF            = "IF"
	ELSE          = "ELSE"
	OR            = "OR"
	AND           = "AND"
	WHILE         = "WHILE"
	FOR           = "FOR"
	FUN           = "FUN"
	RETURN        = "RETURN"
	CLASS         = "CLASS"
	SUPER         = "SUPER"
	THIS          = "THIS"
)

var reservedTokens = map[string]string{
	"(":      LEFT_PAREN,
	")":      RIGHT_PAREN,
	"{":      LEFT_BRACE,
	"}":      RIGHT_BRACE,
	",":      COMMA,
	".":      DOT,
	"-":      MINUS,
	"+":      PLUS,
	";":      SEMICOLON,
	"*":      STAR,
	"=":      EQUAL,
	"==":     EQUAL_EQUAL,
	"!=":     BANG_EQUAL,
	"<":      LESS,
	"<=":     LESS_EQUAL,
	">":      GREATER,
	">=":     GREATER_EQUAL,
	"/":      SLASH,
	"!":      BANG,
	"print":  PRINT,
	"var":    VAR,
	"if":     IF,
	"else":   ELSE,
	"or":     OR,
	"and":    AND,
	"while":  WHILE,
	"for":    FOR,
	"fun":    FUN,
	"return": RETURN,
}

var reservedWords = map[string]string{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}