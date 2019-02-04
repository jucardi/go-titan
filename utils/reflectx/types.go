package reflectx

import (
	"reflect"
)

type ITypeWriter interface {
	SetField(field string, value interface{}) error
	Value() reflect.Value
}

type ITypeReader interface {
	GetFieldHeader() []string
	GetFieldValues() ([]string, error)
}
