package mapper

import (
	"encoding/json"
	"testing"

	"github.com/jucardi/go-testx/testx"
)

const (
	expected1 = `{
  "type": "someType",
  "field1": "someTemplate",
  "field2": 0,
  "some_field": ""
}`
	expected2 = `{
  "type": "someType",
  "field1": "someTemplate",
  "field2": 0,
  "Config": {
    "some_field": "abcde"
  }
}`
)

// BaseConfig contains the standard fields for a formatter configuration
type BaseConfig struct {
	Type   string `json:"type"   yaml:"type"`
	Field1 string `json:"field1" yaml:"field1"`
	Field2 int    `json:"field2" yaml:"field2"`
}

type SomeConfig struct {
	BaseConfig `json:",inline" yaml:",inline"`
	Config     map[string]interface{} `json:",inline" yaml:",inline"`
}

type SomeOtherConfig struct {
	BaseConfig `json:",inline" yaml:",inline"`
	SomeField  string `json:"some_field" yaml:"some_field"`
}

func TestConvert_InlineMapStruct(t *testing.T) {
	testx.Convey("Testing conversion with 2 inline objects in struct", t, func() {
		testx.Convey("From Obj with struct and map inline to struct with extra fields", t, func() {
			obj := &SomeConfig{
				BaseConfig: BaseConfig{
					Type:   "someType",
					Field1: "someTemplate",
				},
				Config: map[string]interface{}{
					"some_field": "abcde",
				},
			}
			target := &SomeOtherConfig{}
			testx.ShouldNotError(Convert(&target, obj, MappingModeJson))
			data, _ := json.MarshalIndent(target, "", "  ")
			testx.ShouldEqual(expected1, string(data))
		})
	})
}

func TestConvert_InlineStructMap(t *testing.T) {
	testx.Convey("Testing conversion with 2 inline objects in struct", t, func() {
		testx.Convey("From Obj with struct and map inline to struct with extra fields", t, func() {
			obj := &SomeOtherConfig{
				BaseConfig: BaseConfig{
					Type:   "someType",
					Field1: "someTemplate",
				},
				SomeField: "abcde",
			}
			target := &SomeConfig{}
			testx.ShouldNotError(Convert(&target, obj))
			data, _ := json.MarshalIndent(target, "", "  ")
			testx.ShouldEqual(expected2, string(data))
		})
	})
}
