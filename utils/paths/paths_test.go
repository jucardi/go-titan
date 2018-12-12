package paths

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombine(t *testing.T) {
	unixPath := Combine("/some/", "/path", "some/", "where", "/in", "some", "machine")
	winPath := Combine("C:\\some\\", "\\path", "some\\", "where", "\\in", "some", "machine")
	httpPath := Combine("http://here.is.some/", "/path", "some/", "where", "/over", "the", "internet")

	assert.Equal(t, unixPath, "/some/path/some/where/in/some/machine")
	assert.Equal(t, winPath, "C:\\some\\path\\some\\where\\in\\some\\machine")
	assert.Equal(t, httpPath, "http://here.is.some/path/some/where/over/the/internet")
}

func TestTempDirAndExists(t *testing.T) {
	tmp, err1 := TempDir()
	exists, err2 := Exists(tmp)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.True(t, exists)
}

func TestDoesNotExists(t *testing.T) {
	exists, err := Exists("/some/very/very/very/random/path/no/one/should/ever/have/in/their/machine")
	assert.Nil(t, err)
	assert.False(t, exists)
}
