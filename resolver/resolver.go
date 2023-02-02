package resolver

import (
	"glox/ast"
	"glox/errors"
	"glox/interpreter"
	"glox/token"
)

type FunctionType int

const (
    NONE FunctionType = iota
    FUNCTION
)

type Resolver struct {
	interp *interpreter.Interpreter
	scopes *Stack[map[string]bool]
    currentFun FunctionType
}

func NewResolver(i *interpreter.Interpreter) *Resolver {
	s := Stack[map[string]bool]{}
	return &Resolver{interp: i, scopes: s.New(), currentFun: NONE}
}

func (r *Resolver) VisitBlockStmt(s *ast.Block) error {
	r.beginScope()
	r.Resolve(s.Statements)
	r.endScope()

	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) Resolve(e interface{}) {
	switch v := e.(type) {
	case []ast.Stmt:
		for _, s := range v {
			r.Resolve(s)
		}

	case ast.Stmt:
		v.Accept(r)

	case ast.Expr:
		v.Accept(r)

	default:
		panic("arg must be of type []ast.Stmt, ast.Stmt or ast.Expr")
	}
}

func (r *Resolver) resolveFunction(f *ast.Function, t FunctionType) error {
    enclosingFun := r.currentFun
    r.currentFun = t

	r.beginScope()
	for _, p := range f.Params {
		r.declare(p)
		r.define(p)
	}

	r.Resolve(f.Body)
	r.endScope()

    r.currentFun = enclosingFun

	return nil
}

func (r *Resolver) VisitVarStmt(s *ast.Var) error {
	r.declare(s.Name)

	if s.Initializer != nil {
		r.Resolve(s.Initializer)
	}

	r.define(s.Name)

	return nil
}

func (r *Resolver) VisitFunctionStmt(s *ast.Function) error {
	r.declare(s.Name)
	r.define(s.Name)

	r.resolveFunction(s, FUNCTION)

	return nil
}

func (r *Resolver) VisitVariableExpr(e *ast.Variable) (interface{}, error) {
	if !r.scopes.IsEmpty() {
        if val, exist := (*r.scopes.Peek())[e.Name.Lexeme()]; exist && !val {
            errors.Error(e.Name, "Can't read local variable in its own initializer.")
        }
	}

	r.resolveLocal(e, e.Name)
	return nil, nil
}

func (r *Resolver) VisitAssignExpr(e *ast.Assign) (interface{}, error) {
	r.Resolve(e.Value)
	r.resolveLocal(e, e.Name)

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(s *ast.Expression) error {
	r.Resolve(s.Exp)

	return nil
}

func (r *Resolver) VisitIfStmt(s *ast.If) error {
	r.Resolve(s.Condition)
	r.Resolve(s.ThenBranch)

	if s.ElseBranch != nil {
		r.Resolve(s.ElseBranch)
	}

	return nil
}

func (r *Resolver) VisitPrintStmt(s *ast.Print) error {
	r.Resolve(s.Exp)
	return nil
}

func (r *Resolver) VisitReturnStmt(s *ast.Return) error {
    if r.currentFun == NONE {
        errors.Error(s.Keyword, "Can't return from top-level code.")
    }

	if s.Value != nil {
		r.Resolve(s.Value)
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(s *ast.While) error {
	r.Resolve(s.Condition)
	r.Resolve(s.Body)

	return nil
}

func (r *Resolver) VisitBinaryExpr(e *ast.Binary) (interface{}, error) {
	r.Resolve(e.Left)
	r.Resolve(e.Right)

	return nil, nil
}

func (r *Resolver) VisitCallExpr(e *ast.Call) (interface{}, error) {
	r.Resolve(e.Callee)

	for _, arg := range e.Arguments {
		r.Resolve(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(e *ast.Grouping) (interface{}, error) {
    r.Resolve(e.Expression)

    return nil, nil
}

func (r *Resolver) VisitLiteralExpr(e *ast.Literal) (interface{}, error) {
    return nil, nil
}

func (r *Resolver) VisitLogicalExpr(e *ast.Logical) (interface{}, error) {
    r.Resolve(e.Left)
    r.Resolve(e.Right)

    return nil, nil
}

func (r *Resolver) VisitUnaryExpr(e *ast.Unary) (interface{}, error) {
    r.Resolve(e.Right)

    return nil, nil
}

func (r *Resolver) resolveLocal(e ast.Expr, name token.Token) {
	for i := r.scopes.Len() - 1; i >= 0; i-- {
		scope := *r.scopes.Get(i)

		if _, ok := scope[name.Lexeme()]; ok {
			r.interp.Resolve(e, r.scopes.Len()-1-i)
		}
	}
}

func (r *Resolver) declare(name token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := *r.scopes.Peek()

    if _, ok := scope[name.Lexeme()]; ok {
        errors.Error(name, "Already a variable with this name in this scope.")
    }

	scope[name.Lexeme()] = false
}

func (r *Resolver) define(name token.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope := *r.scopes.Peek()
	scope[name.Lexeme()] = true
}
