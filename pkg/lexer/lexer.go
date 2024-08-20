package lexer

import (
	"errors"
	"fmt"

	"blorbo/pkg/token"
)

var keywords = map[string]token.TokenType{
	"var":    token.Var,
	"return": token.Return,
	"fn":     token.Fn,
	"struct": token.Struct,
	"for":    token.For,
	"while":  token.While,
	"if":     token.If,
	"else":   token.Else,
	"null":   token.Null,
	"true":   token.True,
	"false":  token.False,
	"and":    token.And,
	"or":     token.Or,
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

type Lexer struct {
	src  string
	pos  int
	line int
}

func New(src string) *Lexer {
	return &Lexer{src: src, line: 1}
}

func (l *Lexer) Scan() ([]token.Token, error) {
	var tokens []token.Token
	var err error

	for l.pos < len(l.src) {
		tok, _err := l.nextToken()
		if _err != nil {
			fmt.Println(_err)
			err = _err
		}

		// Ignore comments
		if tok.Type != token.Comment {
			tokens = append(tokens, tok)
		}
	}

	return tokens, err
}

func (l *Lexer) nextToken() (token.Token, error) {
	c := l.readChar()

	// Consume whitespace
	for ; isWhitespace(c); c = l.readChar() {
		if c == '\n' {
			l.line++
		}
	}

	var tok token.Token

	switch c {
	case 0:
		tok = token.New(token.Eof, "", l.line)
	case '(':
		tok = token.New(token.LeftParen, "(", l.line)
	case ')':
		tok = token.New(token.RightParen, ")", l.line)
	case '{':
		tok = token.New(token.LeftBrace, "{", l.line)
	case '}':
		tok = token.New(token.RightBrace, "}", l.line)
	case '.':
		tok = token.New(token.Dot, ".", l.line)
	case ',':
		tok = token.New(token.Comma, ",", l.line)
	case ';':
		tok = token.New(token.Semicolon, ";", l.line)
	case '*':
		tok = token.New(token.Mul, "*", l.line)
	case '/':
		if l.peekChar() == '/' {
			tok = token.New(token.Comment, l.readComment(), l.line)
		} else {
			tok = token.New(token.Div, "/", l.line)
		}
	case '%':
		tok = token.New(token.Mod, "%", l.line)
	case '+':
		tok = token.New(token.Add, "+", l.line)
	case '-':
		tok = token.New(token.Sub, "-", l.line)
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.Equal, "==", l.line)
		} else {
			tok = token.New(token.Assign, "=", l.line)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.NotEqual, "!=", l.line)
		} else {
			tok = token.New(token.Not, "!", l.line)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.GreaterEqual, ">=", l.line)
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = token.New(token.RightShift, ">>", l.line)
		} else {
			tok = token.New(token.Greater, ">", l.line)
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.New(token.LessEqual, "<=", l.line)
		} else if l.peekChar() == '<' {
			l.readChar()
			tok = token.New(token.LeftShift, "<<", l.line)
		} else {
			tok = token.New(token.Less, "<", l.line)
		}
	case '&':
		tok = token.New(token.BitAnd, "&", l.line)
	case '|':
		tok = token.New(token.BitOr, "|", l.line)
	case '^':
		tok = token.New(token.BitXor, "^", l.line)
	case '~':
		tok = token.New(token.BitNot, "~", l.line)
	case '"':
		str, err := l.readString()
		if err != nil {
			return tok, err
		}
		tok = token.New(token.String, str, l.line)
	default:
		if isAlpha(c) {
			key := l.readIdent()
			val, ok := keywords[key]

			if ok {
				tok = token.New(val, key, l.line)
			} else {
				tok = token.New(token.Ident, key, l.line)
			}
		} else if isDigit(c) {
			tok = token.New(token.Number, l.readNumber(), l.line)
		} else {
			msg := fmt.Sprintf("unexpected character '%c' on line %d", c, l.line)
			return tok, errors.New(msg)
		}
	}

	return tok, nil
}

func (l *Lexer) readChar() byte {
	if l.pos >= len(l.src) {
		return 0
	}

	c := l.src[l.pos]
	l.pos++
	return c
}

func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.src) {
		return 0
	}

	return l.src[l.pos]
}

func (l *Lexer) readNumber() string {
	start := l.pos - 1

	for isDigit(l.peekChar()) {
		l.readChar()
	}

	if l.src[l.pos] == '.' {
		l.readChar()

		for isDigit(l.peekChar()) {
			l.readChar()
		}
	}

	return l.src[start:l.pos]
}

func (l *Lexer) readIdent() string {
	start := l.pos - 1

	for isAlpha(l.peekChar()) || isDigit(l.peekChar()) {
		l.readChar()
	}

	return l.src[start:l.pos]
}

func (l *Lexer) readString() (string, error) {
	start := l.pos
	line := l.line

	for l.peekChar() != '"' && l.peekChar() != 0 {
		if l.peekChar() == '\n' {
			l.line++
		}

		l.readChar()
	}

	if l.peekChar() == 0 {
		msg := fmt.Sprintf("unterminated string on line %d", line)
		return "", errors.New(msg)
	}

	end := l.pos
	l.readChar()

	return l.src[start:end], nil
}

func (l *Lexer) readComment() string {
	start := l.pos

	for l.peekChar() != '\n' && l.peekChar() != 0 {
		l.readChar()
	}

	return l.src[start:l.pos]
}
