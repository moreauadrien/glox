package interpreter

import (
	"glox/ast"
	"glox/environement"
)

type Function struct {
    declaration ast.Function
    closure *environement.Env
}

func NewFunction(declaration ast.Function, closure *environement.Env) Function {
    return Function{declaration: declaration, closure: closure}
}

func (f Function) Call(i *Interpreter, args []interface{}) (interface{}, error) {
    env := environement.NewEnvironement(f.closure)

    for j, p := range f.declaration.Params {
        env.Define(p.Lexeme(), args[j])
    }

    err := i.executeBlock(f.declaration.Body, env)

    if val, ok := err.(Return); ok {
        return val.value, nil
    }

    return nil, err
}

func (f Function) Arity() int {
    return len(f.declaration.Params)
}

func (f Function) String() string {
    return "<fn " + f.declaration.Name.Lexeme() + ">"
}
