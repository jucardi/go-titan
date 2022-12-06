package errors

import (
	"fmt"
	"strings"
)

type defaultBuilder struct {
	flags ErrFlag
	skip  int
}

func (f *defaultBuilder) WithFlags(flags ErrFlag) IErrorBuilder {
	f.skip++
	f.flags = f.flags | flags
	return f
}

func (f *defaultBuilder) New(args ...interface{}) error {
	ret := &stdError{s: fmt.Sprint(args...)}
	ret.InitStack(1 + f.skip)
	ret.f = f.flags
	return ret
}

func (f *defaultBuilder) Format(format string, args ...interface{}) error {
	ret := &stdError{s: fmt.Sprintf(format, args...)}
	ret.InitStack(1 + f.skip)
	ret.f = f.flags
	return ret
}

func (f *defaultBuilder) Join(title string, errs ...string) error {
	if len(errs) == 0 {
		return nil
	}
	strs := append([]string{title}, errs...)
	ret := &stdError{s: strings.Join(strs, "\n - ")}
	ret.InitStack(1 + f.skip)
	ret.f = f.flags
	return ret
}
