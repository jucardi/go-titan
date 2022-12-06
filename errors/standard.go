package errors

var (
	// To ensure *stdError implements the base error definitions defined in this package
	_ IErrorWithStack = (*stdError)(nil)
	_ IErrorWithFlags = (*stdError)(nil)
)

type stdError struct {
	ErrorStackBase
	s string
	f ErrFlag
}

func (e *stdError) Flags() ErrFlag {
	return e.f
}

func (e *stdError) HasFlags(code ErrFlag) bool {
	return HasFlags(e, code)
}

func (e *stdError) Error() string {
	return e.s
}
