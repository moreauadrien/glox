package environement

import (
	"glox/errors"
	"glox/token"
)

type Env struct {
	values    map[string]interface{}
	enclosing *Env
}

func NewEnvironement(enclosing *Env) *Env {
	env := &Env{values: map[string]interface{}{}, enclosing: enclosing}

	return env
}

func (e *Env) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Env) Get(name token.Token) (interface{}, error) {
	if val, ok := e.values[name.Lexeme()]; ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, errors.NewRuntimeErr(name, "Undefined variable '"+name.Lexeme()+"'.")
}

func (e *Env) GetAt(distance int, name string) (interface{}, error) {
	return e.ancestor(distance).values[name], nil
}

func (e *Env) ancestor(distance int) *Env {
    env := e

    for i := 0; i < distance; i++ {
        env = e.enclosing
    }

    return env
}

func (e *Env) Assign(name token.Token, value interface{}) error {
	if _, ok := e.values[name.Lexeme()]; ok {
		e.values[name.Lexeme()] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return errors.NewRuntimeErr(name, "Undefined variable '"+name.Lexeme()+"'.")
}

func (e *Env) AssignAt(distance int, name token.Token, value interface{}) {
    e.ancestor(distance).values[name.Lexeme()] = value
}
