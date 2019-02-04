package reflectx

import (
	"reflect"
	"strings"

	"github.com/jucardi/go-titan/errors"
)

var (
	// Keeps track of types that have been loaded, keeping cache of the field maps and field names
	typesCache = map[string]*typeCache{}
)

type typeCache struct {
	fieldToIndexMap map[string]int
	fieldNames      []string
}

type structWriter struct {
	val          reflect.Value
	fieldToIndex map[string]int
}

func (s *structWriter) SetField(field string, val interface{}) error {
	fIdx, ok := s.fieldToIndex[field]
	if !ok {
		return errors.New("field does not exist")
	}

	fieldVal := s.val.Elem().Field(fIdx)
	fieldType := GetNonPointerType(fieldVal.Type())
	if IsZero(val) {
		return nil
	}

	v, err := GetNonPointerValue(val)
	if err != nil {
		return err
	}
	if v.Type() != fieldType {
		return errors.Format("unable to assign %s to %s", v.Type(), fieldType)
	}
	// Covers any field that are pointers to allow optional (nullable).
	if fieldVal.Kind() == reflect.Ptr {
		ptr := reflect.New(fieldType)
		ptr.Elem().Set(v)
		fieldVal.Set(ptr)
	} else {
		fieldVal.Set(v)
	}

	return nil
}

func (s *structWriter) Value() reflect.Value {
	return s.val
}

// NewStructWriter returns an instance that knows how to write data to a structure.
// If `tag` is provided, the field names can be matched by the value of the provided that
func NewStructWriter(elemType reflect.Type, tag ...string) (ITypeWriter, error) {
	if elemType.Kind() != reflect.Struct {
		return nil, errors.Format("unable to create instance of strucutValue, provided type %s is not a struct", elemType.String())
	}
	t := ""
	if len(tag) > 0 {
		t = tag[0]
	}
	cache, err := registerType(elemType, t)
	if err != nil {
		return nil, errors.Format("unable to generate require mappings for type %s", elemType.String())
	}

	return &structWriter{
		val:          reflect.New(elemType),
		fieldToIndex: cache.fieldToIndexMap,
	}, nil
}

func registerType(elemType reflect.Type, tag string) (*typeCache, error) {
	cacheKey := elemType.String() + "#" + tag
	if c, ok := typesCache[cacheKey]; ok {
		return c, nil
	}

	cache := &typeCache{
		fieldToIndexMap: map[string]int{},
	}

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Name
		if tag != "" {
			if tag, ok := field.Tag.Lookup(tag); ok {
				name = tag
			}
		}
		// To ignore a field when encoding/decoding CSV files
		if name == "-" || strings.ToLower(field.Name[:1]) == field.Name[:1] {
			continue
		}
		if _, ok := cache.fieldToIndexMap[name]; ok {
			return nil, errors.Format("field names in target structure must be uniques (field: %s)", name)
		}
		cache.fieldToIndexMap[name] = i
		cache.fieldNames = append(cache.fieldNames, name)
	}
	typesCache[cacheKey] = cache
	return cache, nil
}
