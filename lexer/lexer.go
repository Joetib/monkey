package lexer

import "monkey/token"

//Lexer : The lexer struct
type Lexer struct {
	input        string
	position     int  // current positino in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination

}

//readChar : Read the current character into ch of lexer
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

}

//NextToken : read and return the next token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)

	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

//peekChar : Helper function to get what the next character might be
//(not the current character we are working with but the next)
// It does not modify the current character
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]

}

//readNumber : reads a number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

//readIdentifier : Read a string
func (l *Lexer) readIdentifier() string {
	position := l.position
	// This method will never be called if the first character
	// was not a letter so we're safe
	// supposedly :-)
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

//readString read and return a string token
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

//skipWhitespace : Helper function to go to next char if current char
// is character without value to us such ar \r,\n, space, \t
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

//isDigit : check if the current character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

//isLetter : check if a given byte input is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

//newToken : Create a new token from a tokentype and a ch : byte
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

//New : construct a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // read the first character
	return l
}
