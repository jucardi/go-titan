package mapper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jucardi/go-titan/errors"
	"github.com/jucardi/go-titan/utils/maps"
	"github.com/jucardi/go-titan/utils/reflectx"
	"gopkg.in/yaml.v3"
)

const (
	MappingModeJson = MappingMode("json")
	MappingModeYaml = MappingMode("yaml")

	jsonTag           = "json"
	yamlTag           = "yaml"
	inlineTagValue    = ",inline"
	omitEmptyTagValye = ",omitempty"

	typeMismatch = "Types Mismatch: "
)

type MappingMode string
type marshalFn func(obj interface{}) ([]byte, error)
type unmarshalFn func(data []byte, obj interface{}) error

type mergeError struct {
	errorType string
	details   string
}

func (e *mergeError) Error() string {
	return e.errorType + e.details
}

// Convert maps the fields of one object to another. Currently it serializes `source` to a YAML or JSON
// and deserializes it again to `target`.
func Convert(target interface{}, source interface{}, mappingMode ...MappingMode) error {
	mode := MappingModeYaml

	if len(mappingMode) > 0 {
		mode = mappingMode[0]
	}

	marshaller, unmarshaller := getMarshallers(mode)
	if marshaller == nil || unmarshaller == nil {
		return fmt.Errorf("invalid mapping mode %s", mode)
	}

	if v, ok := source.(map[interface{}]interface{}); ok {
		source, _ = maps.ConvertMap(v)
	}

	diff, err := getDiff(target, source, mode)
	if err != nil {
		return err
	}

	marshaled, err := marshaller(diff)
	if err != nil {
		return fmt.Errorf("unable to marshal value, %s", err.Error())
	}

	if err := unmarshaller(marshaled, target); err != nil {
		return fmt.Errorf("unable to unmarshal data into the given struct, %s", err.Error())
	}

	appendToInlineMap(target, diff, mode)
	reflectx.Loader().Load(target)

	return nil
}

// Merge attempts to copy the values of each of the provided sources to the provided target.
func Merge(target interface{}, sources ...interface{}) error {
	return merge(MappingModeYaml, target, sources...)
}

func merge(mode MappingMode, target interface{}, sources ...interface{}) error {
	var errs []string
	for _, source := range sources {
		if err := Convert(target, source, mode); err != nil {
			errs = append(errs, err.Error())
		}
	}

	return errors.Join("errors occurred while mering souce values into target", errs...)
}

func getDiff(target interface{}, source interface{}, mode MappingMode) (map[string]interface{}, error) {
	marshaller, unmarshaller := getMarshallers(mode)
	if marshaller == nil || unmarshaller == nil {
		return nil, fmt.Errorf("invalid mapping mode %s", mode)
	}
	sourceM, targetM := map[string]interface{}{}, map[string]interface{}{}

	if data, err := marshaller(source); err != nil {
		return nil, fmt.Errorf("failed to marshal source, %s", err.Error())
	} else if err = unmarshaller(data, &sourceM); err != nil {
		return nil, fmt.Errorf("failed to unmarshal source, %s", err.Error())
	}

	if data, err := marshaller(target); err != nil {
		return nil, fmt.Errorf("failed to marshal target, %s", err.Error())
	} else if err = unmarshaller(data, &targetM); err != nil {
		return nil, fmt.Errorf("failed to unmarshal target, %s", err.Error())
	}

	return getMapDiff(sourceM, targetM)
}

func getMapDiff(src, target map[string]interface{}) (map[string]interface{}, error) {
	ret := map[string]interface{}{}

	for k, v := range src {
		existing, ok := target[k]
		vType, existingType := reflect.TypeOf(v), reflect.TypeOf(existing)
		if vType == nil {
			continue
		}
		if ok && existingType != nil && vType.Kind() != existingType.Kind() && !vType.ConvertibleTo(existingType) {
			return nil, &mergeError{errorType: typeMismatch, details: k}
		} else if reflectx.IsZero(v) {
			continue
		} else if reflect.TypeOf(v).Kind() != reflect.Map {
			if reflectx.IsZero(existing) {
				ret[k] = v
			}
			continue
		}

		nestedSrc, err := maps.ConvertMap(v)
		if err != nil {
			return nil, err
		}
		nestedTarget := map[string]interface{}{}
		if existing != nil {
			nestedTarget, err = maps.ConvertMap(existing)
			if err != nil {
				return nil, err
			}
		}
		diff, err := getMapDiff(nestedSrc, nestedTarget)
		if err != nil {
			return nil, err
		}
		ret[k] = diff
	}
	return ret, nil
}

func getMarshallers(mode MappingMode) (marshalFn, unmarshalFn) {
	switch mode {
	case MappingModeJson:
		return json.Marshal, json.Unmarshal
	case MappingModeYaml:
		return yaml.Marshal, yaml.Unmarshal
	}
	return nil, nil
}

func appendToInlineMap(target interface{}, source map[string]interface{}, mappingMode MappingMode) {
	vTarget, _ := reflectx.GetNonPointerValue(target)

	if source == nil {
		return
	}

	if vTarget.Kind() != reflect.Struct {
		return
	}

	var tag string
	switch mappingMode {
	case MappingModeJson:
		tag = jsonTag
	case MappingModeYaml:
		tag = yamlTag
	}

	for i := 0; i < vTarget.NumField(); i++ {
		fInfo := vTarget.Type().Field(i)
		tagVal := fInfo.Tag.Get(tag)
		key := fInfo.Name
		tag := fInfo.Tag.Get(tag)

		if tag == "-" {
			continue
		}
		if tag != "" {
			key = strings.Split(tagVal, ",")[0]
		}

		fVal, err := reflectx.GetNonPointerValue(vTarget.Field(i))
		if err != nil || fVal.Kind() != reflect.Struct {
			continue
		}

		if fInfo.Anonymous {
			appendToInlineMap(fVal, source, mappingMode)
		} else if nestedSource, err := maps.ConvertMap(source[key]); err == nil {
			appendToInlineMap(fVal, nestedSource, mappingMode)
		}
	}

	valuesSet := map[string]bool{}

	fieldName := getInlineMap(vTarget, tag, valuesSet)
	if fieldName == "" {
		return
	}

	fVal := vTarget.FieldByName(fieldName)

	if !fVal.IsValid() {
		return
	}

	if fVal.IsNil() {
		fVal.Set(reflect.MakeMap(fVal.Type()))
	}

	extraFields := getExtraFieldValues(source, valuesSet, tag)

	for k, v := range extraFields {
		fVal.SetMapIndex(reflect.ValueOf(k), v)
	}
}

func getInlineMap(vTarget reflect.Value, tag string, valuesSet map[string]bool) (ret string) {
	// Find the map field with tag `json:",inline"` or `yaml:",inline"`
	for i := 0; i < vTarget.NumField(); i++ {
		fInfo := vTarget.Type().Field(i)
		fVal := vTarget.Field(i)

		if fInfo.Anonymous && !reflectx.IsNil(fVal) {
			if result := getInlineMap(fVal, tag, valuesSet); result != "" {
				ret = result
			}
			continue
		}

		// Private field case
		if strings.ToLower(fInfo.Name[:1]) == fInfo.Name[:1] {
			continue
		}

		if tagVal := fInfo.Tag.Get(tag); tagVal != inlineTagValue || fInfo.Type.Kind() != reflect.Map {
			switch tagVal {
			case "":
				valuesSet[fInfo.Name] = true
			case "-":
			default:
				tagKey := strings.Split(tagVal, ",")[0]
				valuesSet[tagKey] = true
			}
			continue
		}

		ret = fInfo.Name
	}
	return
}

func getExtraFieldValues(source map[string]interface{}, valuesSet map[string]bool, tag string) map[string]reflect.Value {
	ret := map[string]reflect.Value{}

	if source == nil {
		return ret
	}

	for k, v := range source {
		if _, ok := valuesSet[k]; ok {
			continue
		}
		ret[k] = reflect.ValueOf(v)
	}
	return ret
}
