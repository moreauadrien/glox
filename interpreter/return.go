package interpreter

import "fmt"

type Return struct {
    value interface{}
}

func (r Return) Error() string {
    return fmt.Sprintf("%v", r.value)
}
