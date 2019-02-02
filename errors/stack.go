package errors

import (
	"github.com/jucardi/go-titan/runtime"
)

type IErrorWithStack interface {
	// Stack returns the stack trace of when the error was generated
	Stack() []string
}

type ErrorStackBase struct {
	stack []string
}

// Stack returns the stack trace of when the error was generated
func (e *ErrorStackBase) Stack() []string {
	return e.stack
}

// InitStack initializes the stack trace. If `skip` is provided, it will skip the
// first N entries provided by the `skip` value
func (e *ErrorStackBase) InitStack(skip ...int) {
	s := 1
	if len(skip) > 0 {
		s += skip[0]
	}
	e.stack = runtime.StackLines(s)
}
