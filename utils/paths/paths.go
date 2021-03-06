package paths

import (
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/jucardi/go-titan/logx"
	"github.com/jucardi/go-titan/utils/shutdown"
)

var created []string

func init() {
	shutdown.AddHook(func() error {
		logx.Debug("Temporary directories clean up.")
		for _, v := range created {
			if err := os.RemoveAll(v); err != nil {
				logx.Warnf("Unable to remove temporary directory, %s", v)
			} else {
				logx.Debugf("Temp Dir '%s' deleted.", v)
			}
		}

		return nil
	}, "temp-dirs-cleanup")
}

// Exists indicates whether a file or directory exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// TempDir creates a temporary directory and returns the location
func TempDir() (ret string, err error) {
	ret = path.Join(os.TempDir(), uuid.New().String())
	err = os.Mkdir(ret, 0777)
	addDir(ret)
	return
}

// Combine joins any number of path elements into a single path, adding a
// separating slash if necessary. The result is Cleaned; in particular,
// all empty strings are ignored. Handles resulting double '/'and '\'
func Combine(elem ...string) string {
	p := strings.Replace(path.Join(elem...), ":/", "://", -1)

	if !strings.Contains(p, "\\") {
		return p
	}

	p1 := strings.Replace(p, "/", "\\", -1)
	p2 := strings.Replace(p1, "\\\\", "\\", -1)
	return strings.Replace(p2, "\\\\", "\\", -1)
}

func addDir(dir string) {
	if dir == "" {
		return
	}
	created = append(created, dir)
}
