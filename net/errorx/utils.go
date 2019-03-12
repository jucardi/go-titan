package errorx

import (
	"github.com/jucardi/go-titan/utils/reflectx"
)

// MergeErrors is useful to put together multiple errors in one. Has the following behaviors:
// - Ignores nil entries.
// - If all arguments passed are nil, returns nil.
// - If only one non-nil argument is passed, returns that error.
// - If multiple non-nil arguments are received, creates a new *Error with
//     - Title = "multiple errors occurred"
//     - Message = concatenates all err.Error() separated by `;`
//     - Metadata = an `[]interface{}` instance containing all the non-nil errors received as args. If an err
//       is not a map or struct, it appends the `Error()` result of that error instead to avoid serialization
//       issues.
//
func MergeErrors(err ...error) error {
	var errs []error
	var errType int32

	for _, e := range err {
		if reflectx.IsNil(e) {
			continue
		}
		errs = append(errs, e)
		ex, ok := e.(*Error)
		if !ok {
			ex = wrap(1, e)
		}
		if errType == 0 {
			errType = ex.Code
		} else if errType != ex.Code && errType >= 400 && errType < 500 && ex.Code >= 400 && ex.Code < 500 {
			errType = 500
		} else if ok && ex.Code >= 500 {
			errType = ex.Code
		}
	}

	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	if errType == 0 {
		errType = 500
	}
	return newError(1, 500, "", "", errs...)
}

func MergeErrorx(errs ...*Error) error {
	var converted []error

	for _, err := range errs {
		converted = append(converted, err)
	}
	return MergeErrors(converted...)
}
