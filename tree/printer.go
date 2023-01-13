package tree

/*
import (
	"fmt"
	"strings"
)

type AstPrinter struct {
}

func (p AstPrinter) Print(expr Expr) string {
    s, _ := expr.Accept(p).(string)
    
    return s
}

func (p AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
    return p.parenthesize(expr.Operator.Lexeme(), expr.Left, expr.Right)
}

func (p AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
    return p.parenthesize("group", expr.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
    if expr.Value == nil {
        return nil
    } else {
        return fmt.Sprintf("%v", expr.Value)
    }
}

func (p AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
    return p.parenthesize(expr.Operator.Lexeme(), expr.Right)
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
    b := strings.Builder{}

    b.WriteString("(" + name)

    for _, expr := range exprs {
        b.WriteString(" ")
        s, _ := expr.Accept(p).(string)
        b.WriteString(s)
    }
    b.WriteString(")")

    return b.String()
}
*/
