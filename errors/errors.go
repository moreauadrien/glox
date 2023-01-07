package errors

import "fmt"

var HadError = false

func ErrorAt(line int, message string) {
    report(line, "", message)
}

func report(line int, where string, message string) {
    fmt.Printf("[line %v] Error%v: %v", line, where, message)
    HadError = true
}
