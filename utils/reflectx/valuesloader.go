package reflectx

import (
	"reflect"

	"github.com/jucardi/go-streams/streams"
)

var defaultValuesLoader IValuesLoader

func Loader() IValuesLoader {
	if defaultValuesLoader == nil {
		defaultValuesLoader = newValuesLoader()
	}
	return defaultValuesLoader
}

type IValuesLoader interface {
	Load(obj interface{})
	TagsHandled() []string
}

type loaderHandler func(val string, v reflect.Value, f reflect.StructField)

type tagMapping struct {
	// Tag is the name of the tag that contains the value to assign
	Tag string
	// Loader is a handler that will receive the field object and value so it can set it properly
	Loader loaderHandler
	// Overwrite indicates if any non-zero existing values should be replaced
	Overwrite bool
	// OmitEmpty indicates if empty values should be ignored
	OmitEmpty bool
}

type valuesLoader struct {
	handlers []tagMapping // Using an array instead of a map to allow priorities by order
}

func newValuesLoader() IValuesLoader {
	return &valuesLoader{
		handlers: []tagMapping{
			{Tag: "env", Loader: loadEnvValue, Overwrite: true, OmitEmpty: true},
			{Tag: "envFile", Loader: loadEnvFileValue, Overwrite: true, OmitEmpty: true},
			{Tag: "default", Loader: loadValue, Overwrite: false, OmitEmpty: true},
			{Tag: "defaultFile", Loader: loadFileValue, Overwrite: false, OmitEmpty: true},
		},
	}
}

// Load attempts to loads values into the provided object based on their tags.
func (l *valuesLoader) Load(obj interface{}) {
	val, ok := obj.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(obj)
	}
	t := val.Type()
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		val = val.Elem()
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	if !val.IsValid() {
		return
	}

	for i := 0; t.NumField() > i; i++ {
		v := val.Field(i)
		f := t.Field(i)

		for _, handler := range l.handlers {
			if !handler.Overwrite && !IsZero(v) {
				continue
			}

			if tagValue := f.Tag.Get(handler.Tag); !(tagValue == "-" || (tagValue == "" && handler.OmitEmpty)) {
				handler.Loader(tagValue, v, f)
			}

			fType := f.Type
			for fType.Kind() == reflect.Ptr {
				fType = fType.Elem()
			}
			if fType.Kind() == reflect.Struct {
				l.Load(v)
			}
		}
	}
}

func (l *valuesLoader) TagsHandled() []string {
	return streams.
		From(l.handlers).
		Map(func(i interface{}) interface{} {
			x := i.(tagMapping)
			return x.Tag
		}).
		ToArray().([]string)
}
