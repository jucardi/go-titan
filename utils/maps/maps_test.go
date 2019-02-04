package maps

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"

	"github.com/jucardi/go-testx/assert"
	. "github.com/jucardi/go-testx/testx"
	"github.com/imdario/mergo"
)

var mockString = `{
	"lob": "homes",
	"type": "ms",
	"language": "java",
	"name": "something",
	"db": {
		"enabled": true,
		"type": "mongodb"
	},
	"cache": {
		"enabled": true,
		"type": "redis"
	},
	"null_value": null,
	"array_field": [
		{
			"some_key": "abcdef"
		}, {
			"some_other_key": "xxxxxxxxxxxx"
		}
	]
}`

// FromFile loads the config from a file.  It will exit as fatal if the file does not exist
func getMock() map[string]interface{} {
	ret := make(map[string]interface{})

	if err := json.Unmarshal([]byte(mockString), &ret); err != nil {
		fmt.Printf("Unable to parse JSON. %s", err)
		panic("")
	}
	return ret
}

func TestContains(t *testing.T) {
	mock := getMock()
	assert.True(t, Contains(mock, "language"))
	assert.False(t, Contains(mock, "non-existing-key"))
}

func TestGetValue(t *testing.T) {
	mock := getMock()

	// Test single level of nesting
	val, _ := GetValue(mock, "language")
	assert.Equal(t, "java", val.(string), "Value mismatch")

	// Test get nested value.
	val, _ = GetValue(mock, "db.enabled")
	assert.True(t, val.(bool), "Value mismatch")

	val, _ = GetValue(mock, "array_field[0].some_key")
	assert.Equal(t, "abcdef", val.(string), "Value mismatch")

	val, _ = GetValue(mock, "array_field[1].some_other_key")
	assert.Equal(t, "xxxxxxxxxxxx", val.(string), "Value mismatch")

	_, err := GetValue(mock, "array_field[2].some_other_key")
	assert.Error(t, err)
	assert.Equal(t, "failed to retrieve value at key 'array_field[2].some_other_key'. Index out of range for field 'array_field' (index: 2 | length: 2)", err.Error(), "Unexpected error message")
}

func TestNestedContains(t *testing.T) {
	mock := getMock()
	assert.True(t, Contains(mock, "db.enabled"))
	assert.False(t, Contains(mock, "db.falafel"))
}

func TestGetOrDefault(t *testing.T) {
	mock := getMock()

	assert.Equal(t, GetOrDefault(mock, "cache.type", "something"), "redis")
	assert.Equal(t, GetOrDefault(mock, "some.random.path", "something"), "something")
}

func TestGetNilValue(t *testing.T) {
	mock := getMock()

	v, err := GetValue(mock, "null_value")
	assert.Nil(t, v)
	assert.Nil(t, err)
}

func TestGetNonExistingValue(t *testing.T) {
	mock := getMock()

	v, err := GetValue(mock, "something")
	assert.Nil(t, v)
	assert.Error(t, err)
	assert.Equal(t, "unable to get value by the key 'something'. The value for 'something' is not present", err.Error())
}

func TestSetValue(t *testing.T) {
	mock := getMock()
	_ = SetValue(mock, "db.port", 9999, false)
	val, _ := GetValue(mock, "db.port")
	port := val.(int)
	assert.Equal(t, 9999, port, "Value mismatch")
}

func TestSetValueCreateMissing(t *testing.T) {
	mock := getMock()
	_ = SetValue(mock, "some.nested.value", "value", true)
	val, _ := GetValue(mock, "some.nested.value")
	port := val.(string)
	assert.Equal(t, "value", port, "Value mismatch")
}

func TestSetValueDontCreate(t *testing.T) {
	mock := getMock()
	err := SetValue(mock, "some.nested.value", "value", false)
	assert.NotNil(t, err, "Expected error")
}

func TestConvertMap(t *testing.T) {
	Convey("Testing map conversion", t, func() {
		m := map[interface{}]interface{}{
			"some_key": "some_val",
			"some_nested_obj": map[interface{}]interface{}{
				"field1": 1234,
				"field2": "abcd",
			},
			"some_other_obj": map[string]interface{}{
				"fieldA": 5678,
				"fieldB": "wxyz",
			},
		}
		Convey("When using a map[interface{}]interface{} with other nested map[interface{}]interface{}", t, func() {
			result, err := ConvertMap(m)
			ShouldNotError(err)
			ShouldEqual("map[string]interface {}", reflect.TypeOf(result).String())
			ShouldEqual("map[string]interface {}", reflect.TypeOf(result["some_nested_obj"]).String())
			ShouldEqual("map[string]interface {}", reflect.TypeOf(result["some_other_obj"]).String())
		})
	})
}

func TestStringMapEqual(t *testing.T) {
	Convey("Testing string map comparison", t, func() {
		m1 := map[string]string{
			"test_key1": "test_value1",
			"test_key2": "test_value2",
			"test_key3": "test_value3",
		}

		m2 := map[string]string{
			"test_key3": "test_value3",
			"test_key1": "test_value1",
			"test_key2": "test_value2",
		}

		isEqual := StringMapEqual(m1, m2)
		ShouldBeTrue(isEqual)
	})
}

const (
	yaml1 = `
one:
  one_one: 1
  one_two: true

two:
  two_a:
    two_a_a: a
    two_a_b: b
`

	yaml2 = `
one:
  one_one: one
  one_three: 3

two:
  two_a:
    two_a_b: bee
    two_a_c: c
  two_b:
    two_b_a: 1

three:
  three_one: 123
`
	json1 = `{
  "one": {
    "one_one": 1,
    "one_two": true
  },
  "two": {
    "two_a": {
      "two_a_a": "a",
      "two_a_b": "b"
    }
  }
}`
	json2 = `{
  "one": {
    "one_one": "one",
    "one_three": 3
  },
  "two": {
    "two_a": {
      "two_a_b": "bee",
      "two_a_c": "c"
    },
    "two_b": {
      "two_b_a": 1
    }
  },
  "three": {
    "three_one": 123
  }
}`
)
func TestYamlMerge(t *testing.T) {
	target1 := map[string]interface{}{}
	target2 := map[string]interface{}{}
	_ = yaml.Unmarshal([]byte(yaml1), target1)
	_ = yaml.Unmarshal([]byte(yaml2), target2)
	if err := mergo.Merge(&target1, target2, mergo.WithOverride); err != nil {
		panic(err)
	}

	data, _ := yaml.Marshal(target1)
	println()
	println(string(data))
	println()
	_ = json.Unmarshal([]byte(json1), &target2)
	_ = json.Unmarshal([]byte(json2), &target2)
	data, _ = yaml.Marshal(target1)
	println()
	println(string(data))
	println()
}