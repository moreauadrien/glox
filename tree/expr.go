package tree

import "glox/token"

type Expr interface {
	accept(Visitor) interface{}
}

type Visitor interface {
	visitBinaryExpr(expr *Binary) interface{}
	visitGroupingExpr(expr *Grouping) interface{}
	visitLiteralExpr(expr *Literal) interface{}
	visitUnaryExpr(expr *Unary) interface{}
}

type Binary struct {
	left Expr
	operator token.Token
	right Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) *Binary {
	 return &Binary{left: left, operator: operator, right: right}
}

func (e *Binary) accept(v Visitor) interface{} {
	return v.visitBinaryExpr(e)
}


type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	 return &Grouping{expression: expression}
}

func (e *Grouping) accept(v Visitor) interface{} {
	return v.visitGroupingExpr(e)
}


type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) *Literal {
	 return &Literal{value: value}
}

func (e *Literal) accept(v Visitor) interface{} {
	return v.visitLiteralExpr(e)
}


type Unary struct {
	operator token.Token
	right Expr
}

func NewUnary(operator token.Token, right Expr) *Unary {
	 return &Unary{operator: operator, right: right}
}

func (e *Unary) accept(v Visitor) interface{} {
	return v.visitUnaryExpr(e)
}
