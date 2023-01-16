package interpreter

type Callable interface {
    Call(*Interpreter, []interface{}) (interface{}, error)
    Arity() int
}
