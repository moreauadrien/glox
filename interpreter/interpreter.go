package interpreter

import (
	"fmt"
	"glox/environement"
	"glox/errors"
	"glox/token"
	"glox/ast"
	"reflect"
	"strings"
)

type Interpreter struct {
    env *environement.Env
}

func NewInterpreter() Interpreter {
    return Interpreter{env: environement.NewEnvironement(nil)}
}

func (i Interpreter) Interpret(statements []ast.Stmt) {
    for _, s := range statements {
        if err := i.execute(s); err != nil {
            errors.RuntimeError(err)
            break
        }
    }
}

func (i Interpreter) execute(s ast.Stmt) error {
    return s.Accept(i)
}

func (i Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
    return expr.Accept(i)
}

func (i Interpreter) VisitLiteralExpr(expr *ast.Literal) (interface{}, error) {
    return expr.Value, nil
}

func (i Interpreter) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
    return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
    right, err := i.evaluate(expr.Right)
    if err != nil {
        return nil, err
    }

    switch expr.Operator.Type() {
    case token.BANG:
        return !isTruthy(right), nil

    case token.MINUS:
        if err := checkNumberOperands(expr.Operator, right); err != nil {
            return nil, err
        }

        return -right.(float64), nil
    }

    return nil, nil
}

func (i Interpreter) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
    left, err := i.evaluate(expr.Left)
    if err != nil {
        return nil, err
    }

    right, err := i.evaluate(expr.Right)
    if err != nil {
        return nil, err
    }

    if err := checkNumberOperands(expr.Operator, left, right); err != nil {
        return nil, err
    }

    switch expr.Operator.Type() {
    case token.MINUS:
        return left.(float64) - right.(float64), nil
        
    case token.SLASH:
        if right.(float64) == 0{
            return nil, errors.NewRuntimeErr(expr.Operator, "Divisor must be different from 0")
        }

        return left.(float64) / right.(float64), nil

    case token.STAR:
        return left.(float64) * right.(float64), nil

    case token.PLUS:
        lF, lFOk := left.(float64)
        rF, rFOk := right.(float64)

        if lFOk && rFOk {
            return lF + rF, nil
        }

        _, lSOk := left.(string)
        _, rSOk := right.(string)

        if lSOk || rSOk {
            return stringify(left) + stringify(right), nil
        }

    case token.GREATER:
        return left.(float64) > right.(float64), nil

    case token.GREATER_EQUAL:
        return left.(float64) >= right.(float64), nil

    case token.LESS:
        return left.(float64) < right.(float64), nil

    case token.LESS_EQUAL:
        return left.(float64) <= right.(float64), nil


    case token.EQUAL_EQUAL:
        return isEqual(left, right), nil

    case token.BANG_EQUAL:
        return !isEqual(left, right), nil
    }


    return nil, nil
}

func (i Interpreter) VisitVariableExpr(e *ast.Variable) (interface{}, error) {
    val, err := i.env.Get(e.Name)

    return val, err
}

func (i Interpreter) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
    val, err := i.evaluate(expr.Value)

    if err != nil {
        return nil, err
    }

    i.env.Assign(expr.Name, val)

    return val, nil
}

func (i Interpreter) VisitBlockStmt(s *ast.Block) error {
    return i.executeBlock(s.Statements, environement.NewEnvironement(i.env))
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, env *environement.Env) error {
    prev := i.env

    restoreEnv := func () {
        i.env = prev
    }

    defer restoreEnv()


    i.env = env

    for _, s := range statements {
        err := i.execute(s)
        if err != nil {
            return err
        }

    }

    return nil
}

func (i Interpreter) VisitExpressionStmt(s *ast.Expression) error {
    _, err := i.evaluate(s.Exp)

    return err
}

func (i Interpreter) VisitPrintStmt(s *ast.Print) error {
    val, err := i.evaluate(s.Exp)

    if err != nil {
        return err
    }

    fmt.Println(stringify(val))

    return nil
}


func (i Interpreter) VisitVarStmt(s *ast.Var) error {
    var value interface{}

    if s.Initializer != nil {
        val, err := i.evaluate(s.Initializer)
        
        if err != nil {
            return err
        }

        value = val
    }

    i.env.Define(s.Name.Lexeme(), value)

    return nil
}


func checkNumberOperands(operator token.Token, operands ...interface{}) error {
    switch operator.Type() {
    case token.MINUS, token.SLASH, token.STAR:
        for _, o := range operands {
            if _, ok := o.(float64); !ok {
                plural := ""
                if len(operands) > 1 {
                    plural = "s"
                }

                return errors.NewRuntimeErr(operator, fmt.Sprintf("Operand%v must be number%v.", plural, plural))
            }
        }
    }

    return nil
}

func isTruthy(obj interface{}) bool {
    if obj == nil {
        return false
    }

    if b, ok := obj.(bool); ok {
        return b
    }

    return true 
}

func isEqual(a, b interface{}) bool {
    return reflect.DeepEqual(a, b)
}

func stringify(obj interface{}) string {
    if obj == nil {
        return "nil"
    }


    if _, ok := obj.(float64); ok {
        text := fmt.Sprintf("%v", obj)

        if strings.HasSuffix(text, ".0") {
            text = text[:len(text)-2]
        }

        return text
    }

    return fmt.Sprintf("%v", obj)
}
