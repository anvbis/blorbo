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

	// Bitwise operations
	BitAnd     TokenType = "BitAnd"     // &
	BitOr      TokenType = "BitOr"      // |
	BitXor     TokenType = "BitXor"     // ^
	BitNot     TokenType = "BitNot"     // ~
	RightShift TokenType = "RightShift" // >>
	LeftShift  TokenType = "LeftShift"  // <<

	// Literals
	Ident   TokenType = "Ident"
	Number  TokenType = "Number"
	String  TokenType = "String"
	Comment TokenType = "Comment"

	// Keywords
	Var    TokenType = "Var"    // var
	Return TokenType = "Return" // return
	Fn     TokenType = "Fn"     // fn
	Struct TokenType = "Struct" // struct
	For    TokenType = "For"    // for
	While  TokenType = "While"  // while
	If     TokenType = "If"     // if
	Else   TokenType = "Else"   // else
	Null   TokenType = "Null"   // null
	True   TokenType = "True"   // true
	False  TokenType = "False"  // false
	And    TokenType = "And"    // and
	Or     TokenType = "Or"     // or
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

func New(ty TokenType, literal string, line int) Token {
	return Token{Type: ty, Literal: literal, Line: line}
}
