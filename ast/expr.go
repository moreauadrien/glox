package ast

import "glox/token"

type Expr interface {
	Accept(VisitorExpr) (interface{}, error)
}

type VisitorExpr interface {
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
}

type Assign struct {
	Name token.Token
	Value Expr
}

func NewAssign(Name token.Token, Value Expr) *Assign {
	 return &Assign{Name: Name, Value: Value}
}

func (e *Assign) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitAssignExpr(e)
}


type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(Left Expr, Operator token.Token, Right Expr) *Binary {
	 return &Binary{Left: Left, Operator: Operator, Right: Right}
}

func (e *Binary) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}


type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) *Grouping {
	 return &Grouping{Expression: Expression}
}

func (e *Grouping) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}


type Literal struct {
	Value interface{}
}

func NewLiteral(Value interface{}) *Literal {
	 return &Literal{Value: Value}
}

func (e *Literal) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}


type Unary struct {
	Operator token.Token
	Right Expr
}

func NewUnary(Operator token.Token, Right Expr) *Unary {
	 return &Unary{Operator: Operator, Right: Right}
}

func (e *Unary) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}


type Variable struct {
	Name token.Token
}

func NewVariable(Name token.Token) *Variable {
	 return &Variable{Name: Name}
}

func (e *Variable) Accept(v VisitorExpr) (interface{}, error) {
	return v.VisitVariableExpr(e)
}


