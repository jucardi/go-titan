package reflectx

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/jucardi/go-streams/streams"
	"github.com/jucardi/go-titan/errors"
)

const (
	CloneIgnoreTag = "clone_ignore"
)

// Clone creates a deep clone of the provided object. Will copy references of types that cannot be cloned, otherwise
// will create a full copy of the object. A clone may be returned even if err == nil, for this function, err will provide
// detailed information about the fields that could not be copied.
func Clone(obj interface{}, ignoredFields ...string) (cloned interface{}, err error) {
	val, errs := clone(obj, "", ignoredFields...)
	err = errors.Join("one or more error occurred while cloning the object", errs...)
	if !val.IsZero() {
		cloned = val.Interface()
	}
	return
}

func clone(obj interface{}, path string, ignoredFields ...string) (reflect.Value, []string) {
	val, types, err := getNonPointerValue(obj)
	if err != nil {
		return reflect.Value{}, []string{"failed to get non-pointer-value for '" + path + "' > " + err.Error()}
	}
	var (
		ret  reflect.Value
		errs []string
	)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.String, reflect.Complex64, reflect.Complex128:
		ret = val
	case reflect.Struct:
		ret, errs = cloneStruct(val, path, ignoredFields...)
	case reflect.Slice, reflect.Array:
		ret, errs = cloneSlice(val, path)
	case reflect.Map:
		ret, errs = cloneMap(val, path)
	default:
		return getReflectValue(obj), []string{"field '" + path + "' not cloned. Kind " + val.Kind().String() + " not supported. Reference copied instead"}
	}

	for i := len(types) - 1; i >= 0; i-- {
		switch types[i] {
		case reflect.Ptr:
			ptr := reflect.New(ret.Type())
			ptr.Elem().Set(ret)
			ret = ptr
		}
	}

	return ret, errs
}

func cloneMap(val reflect.Value, path string) (ret reflect.Value, errs []string) {
	if val.IsNil() {
		return reflect.Value{}, nil
	}
	ret = reflect.MakeMap(val.Type())
	iter := val.MapRange()
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		clonedKey, retErrs := clone(k, fmt.Sprintf(path+"[%s].K", k.String()))
		if len(retErrs) > 0 {
			errs = append(errs, retErrs...)
		}
		clonedVal, retErrs := clone(v, fmt.Sprintf(path+"[%s].V", k.String()))
		if len(retErrs) > 0 {
			errs = append(errs, retErrs...)
		}
		ret.SetMapIndex(clonedKey, clonedVal)
	}
	return
}

func cloneSlice(val reflect.Value, path string) (ret reflect.Value, errs []string) {
	if val.IsNil() {
		return reflect.Value{}, nil
	}
	ret = reflect.MakeSlice(val.Type(), val.Len(), val.Cap())
	for i := 0; i < val.Len(); i++ {
		clonedItem, retErrs := clone(val.Index(i), fmt.Sprintf(path+"[%d]", i))
		if len(retErrs) > 0 {
			errs = append(errs, retErrs...)
		}
		if !IsZero(clonedItem) {
			ret.Index(i).Set(clonedItem)
		}
	}
	return
}

func cloneStruct(val reflect.Value, path string, ignoredFields ...string) (reflect.Value, []string) {
	var (
		ret   reflect.Value
		strct reflect.Value
		errs  []string
		t     = val.Type()
	)

	ret = reflect.New(t)
	strct = ret.Elem()

	for i, l := 0, t.NumField(); i < l; i++ {
		f := t.Field(i)
		fPath := path + "." + f.Name

		if _, ignore := f.Tag.Lookup(CloneIgnoreTag); ignore || streams.From(ignoredFields).Contains(fPath[1:]) || unicode.IsLower(rune(f.Name[0])) {
			continue
		}
		fieldVal := val.FieldByName(f.Name)
		if fieldVal.IsZero() {
			continue
		}

		clonedField, retErrs := clone(fieldVal, fPath)
		if len(retErrs) > 0 {
			errs = append(errs, retErrs...)
		}
		if !IsZero(clonedField) {
			strct.FieldByName(f.Name).Set(clonedField)
		}
	}
	return ret.Elem(), errs
}
