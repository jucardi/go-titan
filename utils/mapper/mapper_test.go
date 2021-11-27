package mapper

import (
	"encoding/json"
	"testing"

	"github.com/imdario/mergo"
	. "github.com/jucardi/go-testx/testx"
	"gopkg.in/yaml.v3"
)

type TestStructA struct {
	Field1 string       `json:"field1" yaml:"field1"`
	Field2 int          `json:"field2" yaml:"field2"`
	Field3 *TestStructB `json:"field3" yaml:"field3"`
}

type TestStructB struct {
	Field1 string `json:"field1" yaml:"field1"`
	Field2 int    `json:"field2" yaml:"field2"`
}

var (
	testMap = map[string]interface{}{
		"field1": "abcd",
		"field2": 123,
		"field3": map[string]interface{}{
			"field1": "xyz",
			"field2": 567,
		},
	}
	testStruct = &TestStructA{
		Field1: "abcd",
		Field2: 123,
		Field3: &TestStructB{
			Field1: "xyz",
			Field2: 567,
		},
	}
	mergeObj = map[string]interface{}{
		"extra_field": "something",
	}
)

func TestConvert(t *testing.T) {
	testConvert_MapToStruct(t)
	testConvert_StructToMap(t)
	testConvert_ErrorCases(t)
}

func TestMerge(t *testing.T) {
	testMerge(t)
	testMergeErrorCases(t)
}

func testConvert_MapToStruct(t *testing.T) {
	Convey("Testing map to struct conversion", t, func() {
		target := &TestStructA{}
		Convey("Using mapping mode Mergo", t, func() {
			ShouldNotError(mergo.Map(target, testMap))
			ShouldEqual(testStruct, target)
		})
		Convey("Using mapping mode JSON", t, func() {
			ShouldNotError(Convert(target, testMap, MappingModeJson))
			ShouldEqual(testStruct, target)
		})
		Convey("Using mapping mode YAML", t, func() {
			ShouldNotError(Convert(target, testMap, MappingModeYaml))
			ShouldEqual(testStruct, target)
		})
	})
}

func testConvert_StructToMap(t *testing.T) {
	Convey("Testing struct to map conversion", t, func() {
		Convey("Using mapping mode Mergo", t, func() {
			target := map[string]interface{}{}
			ShouldNotError(mergo.Map(&target, testStruct))
			expected, _ := json.Marshal(testMap)
			actual, _ := json.Marshal(target)
			ShouldEqual(expected, actual)
		})
		Convey("Using mapping mode JSON", t, func() {
			target := map[string]interface{}{}
			ShouldNotError(Convert(&target, testStruct, MappingModeJson))
			expected, _ := json.Marshal(testMap)
			actual, _ := json.Marshal(target)
			ShouldEqual(expected, actual)
		})
		Convey("Using mapping mode YAML", t, func() {
			target := map[string]interface{}{}
			ShouldNotError(Convert(&target, testStruct, MappingModeYaml))
			expected, _ := yaml.Marshal(testMap)
			actual, _ := yaml.Marshal(target)
			ShouldEqual(expected, actual)
		})
	})
}

func testConvert_ErrorCases(t *testing.T) {
	Convey("Testing error cases for coverage", t, func() {
		Convey("Using invalid mapping mode", t, func() {
			target := map[string]interface{}{}
			err := Convert(&target, testStruct, MappingMode("blah"))
			ShouldError(err)
			ShouldEqual("invalid mapping mode blah", err.Error())
		})

		Convey("Not serializable source", t, func() {
			source := map[bool]interface{}{}
			target := map[string]interface{}{}
			err := Convert(&target, source, MappingModeJson)
			ShouldError(err)
		})

		Convey("Invalid target", t, func() {
			err := Convert(nil, testStruct, MappingModeJson)
			ShouldError(err)
		})
	})
}

func testMerge(t *testing.T) {
	Convey("Testing Merge", t, func() {
		target := map[string]interface{}{}
		ShouldNotError(Merge(&target, testStruct, mergeObj))
		expectedMap := cloneMap(testMap)
		for k, v := range mergeObj {
			expectedMap[k] = v
		}
		ShouldEqual(MarshalYAML(expectedMap), MarshalYAML(target))
	})
}

func testMergeErrorCases(t *testing.T) {
	Convey("Testing Merge with one source failing to be serialized", t, func() {
		target := map[string]interface{}{}
		err := merge(MappingModeJson, &target, testStruct, mergeObj, map[bool]interface{}{})
		ShouldError(err)
		expectedMap := cloneMap(testMap)
		for k, v := range mergeObj {
			expectedMap[k] = v
		}
		expected, _ := json.Marshal(expectedMap)
		actual, _ := json.Marshal(target)
		ShouldEqual(expected, actual)
	})
}

func cloneMap(src map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	data, _ := json.Marshal(src)
	_ = json.Unmarshal(data, &ret)
	return ret
}

func TestGetDiff(t *testing.T) {
	Convey("Testing getDiff", t, func() {
		Convey("Using struct as source and map as target", t, func() {
			const expected = `{"field2":123,"field3":{"field2":567}}`
			target := map[string]interface{}{
				"field0": 123,
				"field1": "xxxxx",
				"field3": map[string]interface{}{
					"chaa":   5678,
					"field1": "yyyyyy",
				},
			}
			result, err := getDiff(target, testStruct, MappingModeJson)
			ShouldNotError(err)
			ShouldEqual(expected, MarshalJSON(result))
		})
		Convey("Using map as source and struct as target", t, func() {
			const expected = `{"field2":123,"field3":{"field1":"xyz"}}`
			target := &TestStructA{
				Field1: "xxxxxxx",
				Field3: &TestStructB{
					Field2: 22222222,
				},
			}
			result, err := getDiff(target, testMap, MappingModeJson)
			ShouldNotError(err)
			ShouldEqual(expected, MarshalJSON(result))
		})
	})
}
