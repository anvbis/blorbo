package ast

import (
	"blorbo/pkg/token"
)

type Stmt interface {
}

type BlockStmt struct {
	Body []Stmt
}

type IfStmt struct {
	Cond Expr
	If   Stmt
	Else Stmt
}

type WhileStmt struct {
	Cond Expr
	Body Stmt
}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Inc  Expr
	Body Stmt
}

type FnStmt struct {
	Name   token.Token
	Params []token.Token
	Body   Stmt
}

type VarStmt struct {
	Name  token.Token
	Value Expr
}

type ReturnStmt struct {
	Value Expr
}

type ExprStmt struct {
	Value Expr
}

type Expr interface {
}

type AssignExpr struct {
	Name  token.Token
	Value Expr
}

type BinaryExpr struct {
	Left  Expr
	Op    token.Token
	Right Expr
}

type UnaryExpr struct {
	Op    token.Token
	Right Expr
}

type CallExpr struct {
	Name Expr
	Args []Expr
}

type IdentExpr struct {
	Name token.Token
}

type LiteralExpr struct {
	Value token.Token
}

type Program struct {
	Stmts []Stmt
}
