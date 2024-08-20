package parser

import (
	"errors"
	"fmt"

	"blorbo/pkg/ast"
	"blorbo/pkg/token"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Program -> Stmt* Eof
func (p *Parser) Parse() (*ast.Program, error) {
	program := ast.Program{}

	for !p.checkToken(token.Eof) {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		program.Stmts = append(program.Stmts, stmt)
	}

	return &program, nil
}

// Stmt -> BlockStmt
// | IfStmt
// | FnStmt
// | VarStmt
// | ReturnStmt
// | ExprStmt
func (p *Parser) parseStmt() (ast.Stmt, error) {
	// BlockStmt
	if p.matchToken(token.LeftBrace) {
		return p.parseBlockStmt()
	}

	// IfStmt
	if p.matchToken(token.If) {
		return p.parseIfStmt()
	}

	// ForStmt
	// TODO ForStmt

	// WhileStmt
	// TODO WhileStmt

	// StructStmt
	// TODO StructStmt

	// FnStmt
	if p.matchToken(token.Fn) {
		return p.parseFnStmt()
	}

	// VarStmt
	if p.matchToken(token.Var) {
		return p.parseVarStmt()
	}

	// ReturnStmt
	if p.matchToken(token.Return) {
		return p.parseReturnStmt()
	}

	// ExprStmt
	return p.parseExprStmt()
}

// BlockStmt -> "{" Stmt* "}"
func (p *Parser) parseBlockStmt() (ast.Stmt, error) {
	// Stmt*
	var stmts []ast.Stmt
	for !p.checkToken(token.RightBrace) && !p.checkToken(token.Eof) {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	msg := "expected '}' after statements"
	if _, err := p.expectToken(token.RightBrace, msg); err != nil {
		return nil, err
	}

	return ast.BlockStmt{Body: stmts}, nil
}

// IfStmt -> "if" "(" Expr ")" Stmt ( "else" Stmt )?
func (p *Parser) parseIfStmt() (ast.Stmt, error) {
	msg := "expected '(' after if statement"
	if _, err := p.expectToken(token.LeftParen, msg); err != nil {
		return nil, err
	}

	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	msg = "expected ')' after condition"
	if _, err := p.expectToken(token.RightParen, msg); err != nil {
		return nil, err
	}

	ifStmt, err := p.parseStmt()
	if err != nil {
		return nil, err
	}

	// ( "else" Stmt )?
	var elseStmt ast.Stmt
	if p.matchToken(token.Else) {
		elseStmt, err = p.parseStmt()
		if err != nil {
			return nil, err
		}
	}

	return ast.IfStmt{Cond: cond, If: ifStmt, Else: elseStmt}, nil
}

// FnStmt -> "fn" Ident "(" Params? ")"
// Params -> Ident ( "," Ident )*
func (p *Parser) parseFnStmt() (ast.Stmt, error) {
	msg := "expected function name"
	ident, err := p.expectToken(token.Ident, msg)
	if err != nil {
		return nil, err
	}

	msg = "expected '(' after function name"
	if _, err := p.expectToken(token.LeftParen, msg); err != nil {
		return nil, err
	}

	var params []token.Token
	for ok := true; ok; ok = p.matchToken(token.Comma) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		if expr != nil {
			ident, ok := expr.(ast.IdentExpr)
			if !ok {
				line := p.tokens[p.pos].Line
				msg := fmt.Sprintf("invalid parameter name on line %d", line)
				return nil, errors.New(msg)
			}

			params = append(params, ident.Name)
		}
	}

	msg = "expected ')' after parameters"
	if _, err := p.expectToken(token.RightParen, msg); err != nil {
		return nil, err
	}

	stmt, err := p.parseStmt()
	if err != nil {
		return nil, err
	}

	return ast.FnStmt{Name: ident, Params: params, Body: stmt}, nil
}

// VarStmt -> "var" Ident ( "=" Expr )? ";"
func (p *Parser) parseVarStmt() (ast.Stmt, error) {
	msg := "expected variable name"
	ident, err := p.expectToken(token.Ident, msg)
	if err != nil {
		return nil, err
	}

	// ( "=" Expr)?
	var expr ast.Expr
	if p.matchToken(token.Assign) {
		expr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	msg = "expected ';' after expression"
	if _, err := p.expectToken(token.Semicolon, msg); err != nil {
		return nil, err
	}

	return ast.VarStmt{Name: ident, Value: expr}, nil
}

// ReturnStmt -> "return" Expr ";"
func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	msg := "expected ';' after expression"
	if _, err := p.expectToken(token.Semicolon, msg); err != nil {
		return nil, err
	}

	return ast.ReturnStmt{Value: expr}, nil
}

// ExprStmt -> Expr ";"
func (p *Parser) parseExprStmt() (ast.Stmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	msg := "expected ';' after expression"
	if _, err := p.expectToken(token.Semicolon, msg); err != nil {
		return nil, err
	}

	return expr, nil
}

// Expr -> Assign
func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseAssign()
}

// Assign -> Ident "=" Assign
// | LogicalOr
func (p *Parser) parseAssign() (ast.Expr, error) {
	expr, err := p.parseLogicalOr()
	if err != nil {
		return nil, err
	}

	if p.matchToken(token.Assign) {
		value, err := p.parseAssign()
		if err != nil {
			return nil, err
		}

		expr, ok := expr.(ast.IdentExpr)
		if !ok {
			line := p.tokens[p.pos].Line
			msg := fmt.Sprintf("invalid assignment target on line %d", line)
			return nil, errors.New(msg)
		}

		return ast.AssignExpr{Name: expr.Name, Value: value}, nil
	}

	return expr, nil
}

// LogicalOr -> LogicalAnd ( "or" LogicalAnd )*
func (p *Parser) parseLogicalOr() (ast.Expr, error) {
	expr, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	// ( "or" LogicalAnd )*
	for p.matchToken(token.Or) {
		op := p.prevToken()
		right, err := p.parseLogicalAnd()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// LogicalAnd -> BitwiseOr ( "and" BitwiseOr )*
func (p *Parser) parseLogicalAnd() (ast.Expr, error) {
	expr, err := p.parseBitwiseOr()
	if err != nil {
		return nil, err
	}

	// ( "and" BitwiseOr )*
	for p.matchToken(token.And) {
		op := p.prevToken()
		right, err := p.parseBitwiseOr()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// BitwiseOr -> BitwiseXor ( "|" BitwiseXor )*
func (p *Parser) parseBitwiseOr() (ast.Expr, error) {
	expr, err := p.parseBitwiseXor()
	if err != nil {
		return nil, err
	}

	// ( "|" BitwiseXor )*
	for p.matchToken(token.BitOr) {
		op := p.prevToken()
		right, err := p.parseBitwiseXor()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// BitwiseXor -> BitwiseAnd ( "^" BitwiseAnd )*
func (p *Parser) parseBitwiseXor() (ast.Expr, error) {
	expr, err := p.parseBitwiseAnd()
	if err != nil {
		return nil, err
	}

	// ( "^" BitwiseAnd )*
	for p.matchToken(token.BitXor) {
		op := p.prevToken()
		right, err := p.parseBitwiseAnd()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// BitwiseAnd -> Equality ( "&" Equality )*
func (p *Parser) parseBitwiseAnd() (ast.Expr, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	// ( "&" Equality )*
	for p.matchToken(token.BitAnd) {
		op := p.prevToken()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// Equality -> Comparison ( ( "==" | "!=" ) Comparison )*
func (p *Parser) parseEquality() (ast.Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	// ( ( "==" | "!=" ) Comparison )*
	for p.matchToken(token.Equal) || p.matchToken(token.NotEqual) {
		op := p.prevToken()
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// Comparison -> BitShift ( ( ">" | ">=" | "<" | "<=" ) BitShift )*
func (p *Parser) parseComparison() (ast.Expr, error) {
	expr, err := p.parseBitShift()
	if err != nil {
		return nil, err
	}

	// ( ( ">" | ">=" | "<" | "<=" ) Term )
	for p.matchToken(token.Greater) ||
		p.matchToken(token.GreaterEqual) ||
		p.matchToken(token.Less) ||
		p.matchToken(token.LessEqual) {

		op := p.prevToken()
		right, err := p.parseBitShift()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// BitShift -> Term ( ( ">>" || "<<" ) Term )*
func (p *Parser) parseBitShift() (ast.Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	// ( ( ">>" || "<<" ) Term )*
	for p.matchToken(token.RightShift) || p.matchToken(token.LeftShift) {
		op := p.prevToken()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		
		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// Term -> Factor ( ( "+" | "-" ) Factor )*
func (p *Parser) parseTerm() (ast.Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	// ( ( "+" | "-" ) Term )*
	for p.matchToken(token.Add) || p.matchToken(token.Sub) {
		op := p.prevToken()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// Factor -> Unary ( ( "*" | "/" | "%" ) Unary )*
func (p *Parser) parseFactor() (ast.Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	// ( ( "*" | "/" ) Unary )*
	for p.matchToken(token.Mul) || p.matchToken(token.Div) || p.matchToken(token.Mod) {
		op := p.prevToken()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		expr = ast.BinaryExpr{Left: expr, Op: op, Right: right}
	}

	return expr, nil
}

// Unary -> ( "+" | "-" | "!" | "~" ) Unary
// | Call
func (p *Parser) parseUnary() (ast.Expr, error) {
	// ( "!" | "-" ) Unary
	if p.matchToken(token.Add) ||
		p.matchToken(token.Sub) ||
		p.matchToken(token.Not) ||
		p.matchToken(token.BitNot) {

		op := p.prevToken()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		return ast.UnaryExpr{Op: op, Right: right}, nil
	}

	return p.parseCall()
}

// Call -> Primary ( "(" Args? ")" )*
// Args -> Expression ( "," Expression )*
func (p *Parser) parseCall() (ast.Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.matchToken(token.LeftParen) {
		var args []ast.Expr
		for ok := true; ok; ok = p.matchToken(token.Comma) {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}

			if expr != nil {
				args = append(args, expr)
			}
		}

		msg := "expected ')' after arguments"
		if _, err := p.expectToken(token.RightParen, msg); err != nil {
			return nil, err
		}

		return ast.CallExpr{Name: expr, Args: args}, nil
	}

	return expr, nil
}

// Primary -> Ident
// | Number
// | String
// | "true"
// | "false"
// | "null"
// | "(" Expr ")"
func (p *Parser) parsePrimary() (ast.Expr, error) {
	// Primary -> Ident
	if p.matchToken(token.Ident) {
		return ast.IdentExpr{Name: p.prevToken()}, nil
	}

	// Primary -> Number
	// | String
	// | "true"
	// | "false"
	// | "null"
	if p.matchToken(token.Number) ||
		p.matchToken(token.String) ||
		p.matchToken(token.True) ||
		p.matchToken(token.False) ||
		p.matchToken(token.Null) {

		return ast.LiteralExpr{Value: p.prevToken()}, nil
	}

	// Primary -> "(" Expr ")"
	if p.matchToken(token.LeftParen) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		msg := "expected ')' after expression"
		if _, err := p.expectToken(token.RightParen, msg); err != nil {
			return nil, err
		}

		return expr, nil
	}

	// TODO Unreachable?
	// Or error
	// I think error
	return nil, nil
}

func (p *Parser) checkToken(tok token.TokenType) bool {
	return p.tokens[p.pos].Type == tok
}

func (p *Parser) matchToken(tok token.TokenType) bool {
	if p.checkToken(tok) {
		p.pos++
		return true
	}

	return false
}

func (p *Parser) prevToken() token.Token {
	return p.tokens[p.pos-1]
}

func (p *Parser) expectToken(tok token.TokenType, msg string) (token.Token, error) {
	cur := p.tokens[p.pos]

	if !p.matchToken(tok) {
		msg := fmt.Sprintf("%s on line %d", msg, cur.Line)
		return cur, errors.New(msg)
	}

	return cur, nil
}
