package interpreter

import (
	"fmt"
	"glox/environement"
	"glox/errors"
	"glox/token"
	"glox/tree"
	"reflect"
	"strings"
)

type returnWrapper struct {
    value interface{}
    err error
}

func newReturnWrapper(value interface{}, err error) returnWrapper {
    return returnWrapper{value: value, err: err}
}

type Interpreter struct {
    env *environement.Env
}

func NewInterpreter() Interpreter {
    return Interpreter{env: environement.NewEnvironement(nil)}
}

func (i Interpreter) Interpret(statements []tree.Stmt) {
    for _, s := range statements {
        if err := i.execute(s); err != nil {
            errors.RuntimeError(err)
            break
        }
    }
}

func (i Interpreter) execute(s tree.Stmt) error {
    if err, ok := s.Accept(i).(error); ok {
        return err
    }

    return nil
}

func (i Interpreter) VisitBlockStmt(s *tree.Block) interface{} {
    i.executeBlock(s.Statements, environement.NewEnvironement(i.env))

    return nil
}

func (i *Interpreter) executeBlock(statements []tree.Stmt, env *environement.Env) error {
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

func (i Interpreter) evaluate(expr tree.Expr) interface{} {
    return expr.Accept(i)
}

func (i Interpreter) VisitLiteralExpr(expr *tree.Literal) interface{} {
    return newReturnWrapper(expr.Value, nil)
}

func (i Interpreter) VisitGroupingExpr(expr *tree.Grouping) interface{} {
    return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitUnaryExpr(expr *tree.Unary) interface{} {
    eval := i.evaluate(expr.Right).(returnWrapper)
    if eval.err != nil {
        return eval
    }

    right := eval.value
    

    switch expr.Operator.Type() {
    case token.BANG:
        return newReturnWrapper(!isTruthy(right), nil)

    case token.MINUS:
        if err := checkNumberOperands(expr.Operator, right); err != nil {
            return newReturnWrapper(nil, err)
        }

        return newReturnWrapper(-right.(float64), nil)
    }

    return newReturnWrapper(nil, nil)
}

func (i Interpreter) VisitBinaryExpr(expr *tree.Binary) interface{} {
    evalLeft := i.evaluate(expr.Left).(returnWrapper)
    evalRight := i.evaluate(expr.Right).(returnWrapper)

    if evalLeft.err != nil {
        return evalLeft
    }

    if evalRight.err != nil {
        return evalRight
    }

    left := evalLeft.value
    right := evalRight.value

    if err := checkNumberOperands(expr.Operator, left, right); err != nil {
        return newReturnWrapper(nil, err)
    }

    switch expr.Operator.Type() {
    case token.MINUS:
        return newReturnWrapper(left.(float64) - right.(float64), nil)
        
    case token.SLASH:
        if right.(float64) == 0{
            return newReturnWrapper(nil, errors.NewRuntimeErr(expr.Operator, "Divisor must be different from 0"))
        }

        return newReturnWrapper(left.(float64) / right.(float64), nil)

    case token.STAR:
        return newReturnWrapper(left.(float64) * right.(float64), nil)

    case token.PLUS:
        lF, lFOk := left.(float64)
        rF, rFOk := right.(float64)

        if lFOk && rFOk {
            return newReturnWrapper(lF + rF, nil)
        }

        _, lSOk := left.(string)
        _, rSOk := right.(string)

        if lSOk || rSOk {
            return newReturnWrapper(stringify(left) + stringify(right), nil)
        }

    case token.GREATER:
        return newReturnWrapper(left.(float64) > right.(float64), nil)

    case token.GREATER_EQUAL:
        return newReturnWrapper(left.(float64) >= right.(float64), nil)

    case token.LESS:
        return newReturnWrapper(left.(float64) < right.(float64), nil)

    case token.LESS_EQUAL:
        return newReturnWrapper(left.(float64) <= right.(float64), nil)


    case token.EQUAL_EQUAL:
        return newReturnWrapper(isEqual(left, right), nil)

    case token.BANG_EQUAL:
        return newReturnWrapper(!isEqual(left, right), nil)
    }


    return newReturnWrapper(nil, nil)
}

func (i Interpreter) VisitExpressionStmt(s *tree.Expression) interface{} {
    eval := i.evaluate(s.Exp).(returnWrapper)

    return eval.err
}

func (i Interpreter) VisitPrintStmt(s *tree.Print) interface{} {
    eval := i.evaluate(s.Exp).(returnWrapper)

    if eval.err != nil {
        return eval.err
    }

    fmt.Println(stringify(eval.value))

    return nil
}

func (i Interpreter) VisitVariableExpr(e *tree.Variable) interface{} {
    val, err := i.env.Get(e.Name)

    return newReturnWrapper(val, err)
}

func (i Interpreter) VisitVarStmt(s *tree.Var) interface{} {
    var value interface{}

    if s.Initializer != nil {
        eval := i.evaluate(s.Initializer).(returnWrapper)
        
        if eval.err != nil {
            return eval.err
        }

        value = eval.value
    }

    i.env.Define(s.Name.Lexeme(), value)

    return nil
}

func (i Interpreter) VisitAssignExpr(expr *tree.Assign) interface{} {
    eval := i.evaluate(expr.Value).(returnWrapper)

    if eval.err != nil {
        return eval
    }

    val := eval.value

    i.env.Assign(expr.Name, val)

    return newReturnWrapper(val, nil)
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
