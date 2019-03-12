package info

import (
	"fmt"
	"runtime"
)

var appInfo = map[string]interface{}{
	"version":    Version,
	"built":      Built,
	"processors": runtime.NumCPU(),
	"os_arch":    fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
}

// Add adds the given data object to the provided key. If a value already existed in that key and was not of type `[]interface{}`,
// creates a new `[]interface{}{existing}` and then appends `data` to that array.
func Add(key string, data ...interface{}) {
	if len(data) == 0 {
		return
	}

	if existing, ok := appInfo[key]; ok {
		var list []interface{}

		if l, ok := existing.([]interface{}); ok {
			list = l
		} else {
			list = []interface{}{existing}
		}

		list = append(list, data...)
		appInfo[key] = list
	} else if len(data) > 1 {
		appInfo[key] = data
	} else {
		appInfo[key] = data[0]
	}
}

// Set sets the provided data value to the provided key. Replaces any existing value if the key already exists.
func Set(key string, data interface{}) {
	appInfo[key] = data
}

// Get returns the map containing the info for the application/service
func Get() map[string]interface{} {
	return appInfo
}
