package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals

	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1234
	STRING = "STRING" // eg. "kofi is a boy"

	// OPERATORS

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	LTE      = "<="
	GTE      = ">="
	EQ       = "=="
	NOT_EQ   = "!="
	DOT      = "."

	// DELIMITERS

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	COLON     = ":"

	//KEYWORDS

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	IMPORT   = "IMPORT"
	CLASS    = "CLASS"
	FLOAT    = "FLOAT"
)

//keywords : A map that contains a list of all keywords
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"class":  CLASS,
	"import": IMPORT,
}

//LookupIdent : Checks if an identifier string is a keyword
// if it is a keyword, we return its corresponding token type
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
