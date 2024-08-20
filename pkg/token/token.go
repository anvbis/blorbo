package token

type TokenType string

const (
	// Eof
	Eof TokenType = "Eof"

	// Delimiters
	LeftParen  TokenType = "LeftParen"  // (
	RightParen TokenType = "RightParen" // )
	LeftBrace  TokenType = "LeftBrace"  // {
	RightBrace TokenType = "RightBrace" // }
	Dot        TokenType = "Dot"        // .
	Comma      TokenType = "Comma"      // ,
	Semicolon  TokenType = "Semicolon"  // ;

	// Mathematical operations
	Mul TokenType = "Mul" // *
	Div TokenType = "Div" // /
	Mod TokenType = "Mod" // %
	Add TokenType = "Add" // +
	Sub TokenType = "Sub" // -

	// Logical operations
	Assign       TokenType = "Assign"       // =
	Equal        TokenType = "Equal"        // ==
	Not          TokenType = "Not"          // !
	NotEqual     TokenType = "NotEqual"     // !=
	Greater      TokenType = "Greater"      // >
	GreaterEqual TokenType = "GreaterEqual" // >=
	Less         TokenType = "Less"         // <
	LessEqual    TokenType = "LessEqual"    // <=
	And          TokenType = "And"          // &&
	Or           TokenType = "Or"           // ||

	// Bitwise operations
	BitAnd TokenType = "BitAnd" // &
	BitOr  TokenType = "BitOr"  // |
	BitXor TokenType = "BitXor" // ^

	// Literals
	Ident   = "Ident"
	Number  = "Number"
	String  = "String"
	Comment = "Comment"

	// Keywords
	Var    = "Var"    // var
	Return = "Return" // return
	Fn     = "Fn"     // fn
	Struct = "Struct" // struct
	For    = "For"    // for
	While  = "While"  // while
	If     = "If"     // if
	Else   = "Else"   // else
	Null   = "Null"   // null
	True   = "True"   // true
	False  = "False"  // false
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func New(ty TokenType, literal string, line int) Token {
	return Token{Type: ty, Literal: literal, Line: line}
}
