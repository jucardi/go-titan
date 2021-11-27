package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/jucardi/go-strings/stringx"
	"github.com/jucardi/go-titan/configx"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/utils/osx"
	"github.com/jucardi/go-titan/utils/paths"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

const (
	TestConfigLocation = "fixtures/config-tests.yml"
)

// LoadConfig loads a configuration file meant for tests into the provided `IConfig`.
// If the relative path for the config file is not provided, `LoadConfig` will attempt
// to load the configuration file from the default test config location defined in the
// constant `TestConfigLocation`
func LoadConfig(relativePath ...string) {
	file := stringx.GetOrDefault(TestConfigLocation, relativePath...)
	if configx.FromFile(file) != nil {
		panic("failed to load test configuration")
	}
}

// DefaultTestConfigPath returns the default test config file location relative to the
// project location.
func DefaultTestConfigPath() string {
	return AbsolutePath(TestConfigLocation)
}

// AbsolutePath will combine the project absolute path with the relative
// file/path starting from the root.
//
// For example, if the argument of "test_assets/somefile.json" is specified the result
// returned would be the [ Project Root Absolute Path ] + test_assets/somefile.json
func AbsolutePath(file string) string {
	return paths.Combine(osx.ProjectRoot(), file)
}

// FileToReader loads JSON assets to a reader for tests with request bodies.
func FileToReader(filename string) io.Reader {
	if content, err := ioutil.ReadFile(filename); err != nil {
		logx.Fatalf("Unable to read file '%s': %s", filename, err.Error())
	} else {
		return bytes.NewReader(content)
	}
	return nil
}

// FileToObj loads a JSON asset to an object.
func FileToObj(filename string, target interface{}) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		logx.Fatalf("Unable to read file '%s': %s", filename, err.Error())
	}

	switch path.Ext(strings.ToLower(filename)) {
	case ".yml":
		fallthrough
	case ".yaml":
		if err := yaml.Unmarshal(content, target); err != nil {
			logx.Fatalf("Unable to parse YAML file '%s': %s", filename, err)
		}
	case ".json":
		fallthrough
	default:
		if err := json.Unmarshal(content, target); err != nil {
			logx.Fatalf("Unable to parse JSON file '%s': %s", filename, err)
		}
	}
}

// FileToProtoObj loads a JSON asset into a protobuf object
func FileToProtoObj(filename string, target proto.Message) {
	if content, err := ioutil.ReadFile(filename); err != nil {
		logx.Fatalf("Unable to read file '%s': %s", filename, err.Error())
	} else if err = protojson.Unmarshal(content, target); err != nil {
		logx.Fatalf("Unable to parse JSON  file '%s': %s", filename, err)
	}
}

// FileToProtoBytes loads a JSON asset into bytes for the specified protocol buffer
func FileToProtoBytes(filename string, target proto.Message) []byte {
	FileToProtoObj(filename, target)
	if data, err := proto.Marshal(target); err != nil {
		logx.Fatalf("Unable to marshal proto to bytes", err)
	} else {
		return data
	}
	return nil
}

// ObjToReader loads an object to an io.Reader, useful to load request, or response bodies for mocks.
func ObjToReader(obj interface{}) io.Reader {
	if content, err := json.MarshalIndent(obj, "", "  "); err != nil {
		logx.Fatalf("Unable to marshal object: %s", err.Error())
	} else {
		return bytes.NewReader(content)
	}
	return nil
}
