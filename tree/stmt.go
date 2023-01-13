package tree

import "glox/token"

type Stmt interface {
	Accept(VisitorStmt) interface{}
}

type VisitorStmt interface {
	VisitBlockStmt(stmt *Block) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	VisitPrintStmt(stmt *Print) interface{}
	VisitVarStmt(stmt *Var) interface{}
}

type Block struct {
	Statements []Stmt
}

func NewBlock(Statements []Stmt) *Block {
	 return &Block{Statements: Statements}
}

func (e *Block) Accept(v VisitorStmt) interface{} {
	return v.VisitBlockStmt(e)
}


type Expression struct {
	Exp Expr
}

func NewExpression(Exp Expr) *Expression {
	 return &Expression{Exp: Exp}
}

func (e *Expression) Accept(v VisitorStmt) interface{} {
	return v.VisitExpressionStmt(e)
}


type Print struct {
	Exp Expr
}

func NewPrint(Exp Expr) *Print {
	 return &Print{Exp: Exp}
}

func (e *Print) Accept(v VisitorStmt) interface{} {
	return v.VisitPrintStmt(e)
}


type Var struct {
	Name token.Token
	Initializer Expr
}

func NewVar(Name token.Token, Initializer Expr) *Var {
	 return &Var{Name: Name, Initializer: Initializer}
}

func (e *Var) Accept(v VisitorStmt) interface{} {
	return v.VisitVarStmt(e)
}


