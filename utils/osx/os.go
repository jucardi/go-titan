package osx

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/jucardi/go-titan/utils/paths"
)

var (
	projectRootRE = regexp.MustCompile(".*/(src/\\S+?/\\S+?/\\S+?)+/")
	projectRoot   string
)

type fileWriter struct {
	filename string
	file     *os.File
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *fileWriter) Close() error {
	return w.file.Close()
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string, overwrite bool, perm ...os.FileMode) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return err
	}
	if _, err := os.Stat(dst); !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("destination file '%s' already exists, overwrite flag set to false, %s", src, err.Error())
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	p := os.FileMode(0655)
	if len(perm) > 0 {
		p = perm[0]
	}
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, p)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// ProjectRoot is the root directory of the calling project
func ProjectRoot() string {
	if projectRoot != "" {
		return projectRoot
	}

	workdir, err := os.Getwd()
	if err != nil {
		workdir = os.Getenv("PWD")
	}

	if strings.Contains(workdir, ":\\") {
		workdir = strings.Replace(workdir, "\\", "/", -1)
	}

	srcMatch := projectRootRE.FindStringSubmatch(workdir)
	if len(srcMatch) > 0 {
		projectRoot = paths.Combine(os.Getenv("GOPATH"), srcMatch[1])
	}

	return projectRoot
}

// NewFileWriter creates a new instance of a file writer
func NewFileWriter(filename string) (io.WriteCloser, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &fileWriter{filename: filename, file: f}, nil
}
