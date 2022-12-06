package errors

type ErrFlag uint

func (e ErrFlag) String() string {
	if str, ok := FlagsToStringMap[e]; ok {
		return str
	}
	return "unknown"
}

// The following is a list of predefined flags for common errors, however, a custom flag list may be
// defined and used instead of this one.
//
// If a custom list of flags is used, a new map[ErrFlag]string should be assigned to replace the
// default map `FlagsToStringMap` for a proper string conversion.
//
// If also using the `net` and/or the `net/errorx` package, a new ErrFlag to HTTP Status map should
// be assigned to the variable `ErrFlagsMapping` in the `net/errorx` package for proper error wrapping
const (
	FlagUnhandled ErrFlag = 1 << iota
	FlagBadRequest
	FlagNotFound
	FlagUnauthorized
	FlagNotImplemented
	FlagOperationTimeout
	FlagConflict
)

var (
	// FlagsToStringMap is a map of ErrFlags to a friendly string representation. As mentioned above,
	// this map should be replaced if using a different set of error flags than then ones defined in
	// this package, as well as the mapping of ErrFlags to HTTP status codes in `net/errorx`
	FlagsToStringMap = map[ErrFlag]string{
		FlagUnhandled:        "unhandled",
		FlagBadRequest:       "bad request",
		FlagNotFound:         "not found",
		FlagUnauthorized:     "unauthorized",
		FlagNotImplemented:   "not implemented",
		FlagOperationTimeout: "operation timeout",
		FlagConflict:         "conflict",
	}
)

// IErrorWithStack defines an error type that supports returning a stack trace
type IErrorWithStack interface {
	error

	// Stack returns the stack trace of when the error was generated
	Stack() []string
}

// IErrorWithFlags defines an error type that supports returning a common error code.
// Error codes are meant to be binary flags so multiple codes can be appended to the same error instance
type IErrorWithFlags interface {
	error

	// Flags returns the code flags contained by this error
	Flags() ErrFlag

	// HasFlags indicates whether this error instance has the provided flag(s)
	HasFlags(code ErrFlag) bool
}

type IErrorBuilder interface {
	// New is equivalent `errors.New` from the "errors" package, with the difference that allows
	// multiple arguments to be passed which will be processed by fmt.Sprint
	//
	// Appends the stack trace data which can be retrieved with `Stack()`
	New(args ...interface{}) error

	// Format is equivalent to `fmt.Errorf`
	//
	// Appends the stack trace data which can be retrieved with `Stack()`
	Format(format string, args ...interface{}) error

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
	Join(title string, errs ...string) error

	// WithFlags sets the provided flags into the current builder, keeping all previously existing flags
	WithFlags(flags ErrFlag) IErrorBuilder
}
