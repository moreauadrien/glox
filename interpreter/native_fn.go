package interpreter

type callFunction func(*Interpreter, []interface{}) (interface{}, error)

type NativeFunction struct {
    arity int
    call callFunction
}

func (n NativeFunction) Arity() int {
    return n.arity
}

func (n NativeFunction) Call(i *Interpreter, args []interface{}) (interface{}, error) {
    return n.call(i, args)
}

func (n NativeFunction) String() string {
    return "<native fn>"
}
