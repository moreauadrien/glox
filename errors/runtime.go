package errors

import (
	"fmt"
	"glox/token"
)

type RuntimeErr struct {
    token token.Token
    message string
}

func NewRuntimeErr(t token.Token, message string) RuntimeErr {
    return RuntimeErr{token: t, message: message}
}

func (e RuntimeErr) Error() string {
    return fmt.Sprintf("%v\n[line %v]", e.message, e.token.Line())
}
