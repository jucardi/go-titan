package reflectx

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

var (
	settersByKind = map[reflect.Kind]func(v reflect.Value, val []byte){
		reflect.String:  loadString,
		reflect.Int:     loadInt,
		reflect.Int8:    loadInt8,
		reflect.Int16:   loadInt16,
		reflect.Int32:   loadInt32,
		reflect.Int64:   loadInt64,
		reflect.Bool:    loadBool,
		reflect.Float64: loadFloat64,
		reflect.Float32: loadFloat32,
	}
	settersByType = map[string]func(v reflect.Value, val []byte){
		"[]byte": loadBytes,
	}
)

func loadEnvValue(val string, v reflect.Value, f reflect.StructField) {
	loadValue(os.Getenv(val), v, f)
}

func loadEnvFileValue(val string, v reflect.Value, f reflect.StructField) {
	loadFileValue(os.Getenv(val), v, f)
}

func loadValue(val string, v reflect.Value, f reflect.StructField) {
	if val == "" {
		return
	}
	if setter, ok := settersByKind[v.Kind()]; ok {
		setter(v, []byte(val))
	} else if setter, ok = settersByType[v.Type().String()]; ok {
		setter(v, []byte(val))
	} else {
		log().Warnf("loader not found for kind:%s or type:%s (field:%s)", v.Kind(), v.Type(), f.Name)
	}
}

func loadFileValue(file string, v reflect.Value, f reflect.StructField) {
	if file == "" {
		return
	}
	file = os.ExpandEnv(file)
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log().Warnf("failed to load file contents for config field %s, %s", f.Name, err.Error())
		return
	}
	// TODO: add json and yaml deserialization based of file extension
	if setter, ok := settersByKind[v.Kind()]; ok {
		setter(v, contents)
	} else if setter, ok = settersByType[v.Type().String()]; ok {
		setter(v, contents)
	} else {
		log().Warnf("loader not found for kind:%s or type:%s (field:%s)", v.Kind(), v.Type(), f.Name)
	}
}

func loadString(v reflect.Value, val []byte) {
	v.Set(reflect.ValueOf(string(val)).Convert(v.Type()))
}

func loadInt8(v reflect.Value, val []byte) {
	intVal, err := strconv.ParseInt(string(val), 10, 8)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not int8, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(int8(intVal)).Convert(v.Type()))
}

func loadInt16(v reflect.Value, val []byte) {
	intVal, err := strconv.ParseInt(string(val), 10, 16)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not int16, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(int16(intVal)).Convert(v.Type()))
}

func loadInt32(v reflect.Value, val []byte) {
	intVal, err := strconv.ParseInt(string(val), 10, 32)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not int32, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(int32(intVal)).Convert(v.Type()))
}

func loadInt64(v reflect.Value, val []byte) {
	intVal, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not int64, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(intVal).Convert(v.Type()))
}

func loadInt(v reflect.Value, val []byte) {
	intVal, err := strconv.Atoi(string(val))
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not int, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(intVal).Convert(v.Type()))
}

func loadFloat32(v reflect.Value, val []byte) {
	floatVal, err := strconv.ParseFloat(string(val), 32)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not float32, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(float32(floatVal)).Convert(v.Type()))
}

func loadFloat64(v reflect.Value, val []byte) {
	floatVal, err := strconv.ParseFloat(string(val), 64)
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not float64, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(floatVal).Convert(v.Type()))
}

func loadBool(v reflect.Value, val []byte) {
	bVal, err := strconv.ParseBool(string(val))
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not bool, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(bVal).Convert(v.Type()))
}

func loadBytes(v reflect.Value, val []byte) {
	bytesVal, err := base64.StdEncoding.DecodeString(string(val))
	if err != nil {
		log().Errorf("failed to load env variable, '%s' is not a base64 encoded string, %s", val, err.Error())
	}

	v.Set(reflect.ValueOf(bytesVal))
}
