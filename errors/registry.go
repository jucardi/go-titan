package errors

import (
	"reflect"
	"sync"

	"github.com/jucardi/go-titan/utils/reflectx"
)

var (
	typeRegistry = map[reflect.Type]ErrFlag{}
	errRegistry  = map[error]ErrFlag{}

	typeMux = &sync.RWMutex{}
	errMux  = &sync.RWMutex{}
)

// GetFlagsByType attempts to obtain the error flags registered to an error type.
// Returns 0 if none found or if failed to obtain the error type
func GetFlagsByType(err error) ErrFlag {
	t := reflectx.GetNonPointerType(err)
	if t.Kind() == reflect.Invalid {
		return 0
	}

	typeMux.Lock()
	defer typeMux.Unlock()

	if c, ok := typeRegistry[t]; ok {
		return c
	}
	return 0
}

// GetFlagsByInstance attempts to obtain the error flags registered to an error instance.
// Returns 0 if none found or if err == nil
func GetFlagsByInstance(err error) ErrFlag {
	if reflectx.IsNil(err) {
		return 0
	}

	errMux.Lock()
	defer errMux.Unlock()

	if c, ok := errRegistry[err]; ok {
		return c
	}
	return 0
}

// SetFlagsByType registers an error type into the registry to be tied to a provided error flag.
// If the error type was previously registered, it will set the provided flag(s) to the existing
// value instead.
func SetFlagsByType(err error, flags ErrFlag) error {
	t := reflectx.GetNonPointerType(err)
	if t.Kind() == reflect.Invalid {
		return New("failed to register the provided error, unable to determine its type")
	}

	typeMux.Lock()
	defer typeMux.Unlock()

	if c, ok := typeRegistry[t]; ok {
		typeRegistry[t] = c | flags
	} else {
		typeRegistry[t] = flags
	}
	return nil
}

// SetFlagsByInstance registers an error instance into the registry to be tied to a provided error flag.
// If the error instance was previously registered, it will set the provided flag(s) to the existing
// value instead.
func SetFlagsByInstance(err error, flags ErrFlag) {
	if reflectx.IsNil(err) {
		return
	}

	errMux.Lock()
	defer errMux.Unlock()

	if c, ok := errRegistry[err]; ok {
		errRegistry[err] = c | flags
	} else {
		errRegistry[err] = flags
	}
}

// UnsetFlagsByType unsets the provided flags from the existing value tied to the specified type.
// Does nothing if the error instance was not previously registered.
func UnsetFlagsByType(err error, flags ErrFlag) error {
	t := reflectx.GetNonPointerType(err)
	if t.Kind() == reflect.Invalid {
		return New("failed to unset the flags for the provided error, unable to determine its type")
	}

	typeMux.Lock()
	defer typeMux.Unlock()

	if c, ok := typeRegistry[t]; ok && c&flags == flags {
		typeRegistry[t] = c ^ flags
	}
	return nil
}

// UnsetFlagsByInstance unsets the provided flags from the existing value tied to the specified instance.
// Does nothing if the error instance was not previously registered.
func UnsetFlagsByInstance(err error, flags ErrFlag) {
	if reflectx.IsNil(err) {
		return
	}

	errMux.Lock()
	defer errMux.Unlock()

	if c, ok := errRegistry[err]; ok && c&flags == flags {
		errRegistry[err] = c ^ flags
	}
}

// DeregisterType attempts to remove any flags assigned to a previously registered error type. Does
// nothing if the type was not previously registered
func DeregisterType(err error) error {
	t := reflectx.GetNonPointerType(err)
	if t.Kind() == reflect.Invalid {
		return New("failed to deregister the provided error, unable to determine its type")
	}

	typeMux.Lock()
	defer typeMux.Unlock()

	if _, ok := typeRegistry[t]; !ok {
		return nil
	}
	delete(typeRegistry, t)
	return nil
}

// DeregisterInstance attempts to remove any flags assigned to a previously registered error instance.
// Does nothing if the instance was not previously registered
func DeregisterInstance(err error) {
	if reflectx.IsNil(err) {
		return
	}

	errMux.Lock()
	defer errMux.Unlock()

	if _, ok := errRegistry[err]; !ok {
		return
	}
	delete(errRegistry, err)
}

// GetFlags attempts to calculate any flags associated to the provided error. This is determined by
// an OR bitwise operation of the following values:
//
//   - If the error implements IErrorWithFlags, the return value of e.Flags()
//   - Flags registered by instance (if any)
//   - Flags registered by error type (if any)
//
// Returns 0 if err == nil
//
func GetFlags(err error) ErrFlag {
	if reflectx.IsNil(err) {
		return 0
	}

	var ret ErrFlag

	if withCode, ok := err.(IErrorWithFlags); ok {
		ret = ret | withCode.Flags()
	}

	ret = ret | GetFlagsByInstance(err)
	ret = ret | GetFlagsByType(err)

	return ret
}

// Match attempts to determine if a provided error contains the provided flag(s). This is determined
// by any of the following scenarios:
//
//   - The provided error instance implements IErrorWithFlags and e.HasFlags(err, flags) returns true.
//   - The provided error instance is registered and matches the provided code
//   - The provided error type is registered and matches the provided code
//
func Match(err error, code ErrFlag) bool {
	if reflectx.IsNil(err) {
		return false
	}

	if withCode, ok := err.(IErrorWithFlags); ok && withCode.HasFlags(code) {
		return true
	}
	return isInstanceMatch(err, code) || isTypeMatch(err, code)
}

func isInstanceMatch(err error, code ErrFlag) bool {
	errMux.RLock()
	defer errMux.RUnlock()

	if c, ok := errRegistry[err]; ok {
		return c&code == code
	}
	return false
}

func isTypeMatch(err error, code ErrFlag) bool {
	typeMux.RLock()
	defer typeMux.RUnlock()

	if t := reflectx.GetNonPointerType(err); t.Kind() == reflect.Invalid {
		return false
	} else if c, ok := typeRegistry[t]; ok {
		return c&code == code
	}
	return false
}
