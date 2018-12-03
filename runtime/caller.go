package runtime

import (
	"fmt"
	"runtime"
	"strings"
)

// Caller returns the caller file and line using the `runtine` package.
// Removes the $GOPATH/src prefix from the result.
// Returns "" if it was unable to fetch the caller data.
//
// NOTE: For the successful removal of the prefix, requires the solution to be compiled inside the `GOPATH`
//       Returns the filename if it was unsuccessful removing the prefix.
//
func Caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	return getCallerString(file) + ":" + fmt.Sprint(line)
}

func getCallerString(file string) string {
	file = strings.Replace(file, "\\", "/", -1)
	split := strings.Split(file, "/")

	var foundSrc bool
	var ret []string

	for _, segment := range split {
		if foundSrc {
			ret = append(ret, segment)
			continue
		}

		if strings.ToLower(segment) == "src" {
			foundSrc = true
		}
	}

	if len(ret) == 0 {
		return file
	}
	return strings.Join(ret, "/")
}
