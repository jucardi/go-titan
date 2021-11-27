package errors

import (
	"fmt"
	"strings"
)

type strError struct {
	ErrorStackBase
	s string
}

// Error returns the error message
func (e *strError) Error() string {
	return e.s
}

// New is equivalent `errors.New` from the "errors" package, with the difference that allows
// multiple arguments to be passed which will be processed by fmt.Sprint
//
// Appends the stack trace data which can be retrieved with `Stack()`
func New(args ...interface{}) error {
	ret := &strError{s: fmt.Sprint(args...)}
	ret.InitStack(1)
	return ret
}

// Format is equivalent to `fmt.Errorf`
//
// Appends the stack trace data which can be retrieved with `Stack()`
func Format(format string, args ...interface{}) error {
	ret := &strError{s: fmt.Sprintf(format, args...)}
	ret.InitStack(1)
	return ret
}

// Join creates a new error from a provided list of error strings and a title.
// Returns nil if the list is nil or empty.
//
// Eg:
//   errs := []string{"file not found", "failed to communicate with db"}
//   err := errors.Join("Failed finish transaction", errs...)
//   println(err.Error())
//
// Output:
//  > Failed finish transaction
//  >  - file not found
//  >  - failed to communicate with db
//
// Appends the caller data which can be retrieved with `Caller()`
func Join(title string, errs ...string) error {
	if len(errs) == 0 {
		return nil
	}
	strs := append([]string{title}, errs...)
	ret := &strError{s: strings.Join(strs, "\n - ")}
	ret.InitStack(1)
	return ret
}
