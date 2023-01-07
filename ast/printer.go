package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (p AstPrinter) print(expr Expr) string {
    s, _ := expr.accept(p).(string)
    
    return s
}

func (p AstPrinter) visitBinaryExpr(expr *Binary) interface{} {
    return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p AstPrinter) visitGroupingExpr(expr *Grouping) interface{} {
    return p.parenthesize("group", expr.expression)
}

func (p AstPrinter) visitLiteralExpr(expr *Literal) interface{} {
    if expr.value == nil {
        return nil
    } else {
        return fmt.Sprintf("%v", expr.value)
    }
}

func (p AstPrinter) visitUnaryExpr(expr *Unary) interface{} {
    return p.parenthesize(expr.operator.Lexeme, expr.right)
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
    b := strings.Builder{}

    b.WriteString("(" + name)

    for _, expr := range exprs {
        b.WriteString(" ")
        s, _ := expr.accept(p).(string)
        b.WriteString(s)
    }
    b.WriteString(")")

    return b.String()
}
