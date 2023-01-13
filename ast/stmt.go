package ast

import "glox/token"

type Stmt interface {
	Accept(VisitorStmt) error
}

type VisitorStmt interface {
	VisitBlockStmt(stmt *Block) error
	VisitExpressionStmt(stmt *Expression) error
	VisitPrintStmt(stmt *Print) error
	VisitVarStmt(stmt *Var) error
}

type Block struct {
	Statements []Stmt
}

func NewBlock(Statements []Stmt) *Block {
	 return &Block{Statements: Statements}
}

func (e *Block) Accept(v VisitorStmt) error {
	return v.VisitBlockStmt(e)
}


type Expression struct {
	Exp Expr
}

func NewExpression(Exp Expr) *Expression {
	 return &Expression{Exp: Exp}
}

func (e *Expression) Accept(v VisitorStmt) error {
	return v.VisitExpressionStmt(e)
}


type Print struct {
	Exp Expr
}

func NewPrint(Exp Expr) *Print {
	 return &Print{Exp: Exp}
}

func (e *Print) Accept(v VisitorStmt) error {
	return v.VisitPrintStmt(e)
}


type Var struct {
	Name token.Token
	Initializer Expr
}

func NewVar(Name token.Token, Initializer Expr) *Var {
	 return &Var{Name: Name, Initializer: Initializer}
}

func (e *Var) Accept(v VisitorStmt) error {
	return v.VisitVarStmt(e)
}


