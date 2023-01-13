package errors

import (
	"fmt"
	"glox/token"
)

var HadError = false
var HadRuntimeError = false

func ErrorAt(line int, message string) {
    report(line, "", message)
}

func report(line int, where string, message string) {
    fmt.Printf("[line %v] Error%v: %v\n", line, where, message)
    HadError = true
}

func Error(t token.Token, message string) {
    if t.Type() == token.EOF {
        report(t.Line(), " at end", message)
    } else {
        report(t.Line(), " at '" + t.Lexeme() + "'", message)
    }
}

func RuntimeError(err error) {
    fmt.Println(err.Error())
    HadRuntimeError = true
}
