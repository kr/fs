package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem defines the methods of an abstract filesystem.
type FileSystem interface {

	// ReadDir reads the directory named by dirname and returns a
	// list of directory entries.
	ReadDir(dirname string) ([]os.FileInfo, error)

	// Lstat returns a FileInfo describing the named file. If the file is a
	// symbolic link, the returned FileInfo describes the symbolic link. Lstat
	// makes no attempt to follow the link.
	Lstat(name string) (os.FileInfo, error)

	// PathSeparator returns the FileSystem specific path separator.
	PathSeparator() byte
}

// fs represents a FileSystem provided by the os package.
type fs struct{}

func (f *fs) ReadDir(dirname string) ([]os.FileInfo, error) { return ioutil.ReadDir(dirname) }

func (f *fs) Lstat(name string) (os.FileInfo, error) { return os.Lstat(name) }

func (f *fs) PathSeparator() byte { return os.PathSeparator }

// Join joins any number of path elements into a single path, adding
// a Separator if necessary. The result is Cleaned, in particular
// all empty strings are ignored.
func Join(fs FileSystem, elem ...string) string {
	sep := string(fs.PathSeparator())
	for i, e := range elem {
		if e != "" {
			return filepath.Clean(strings.Join(elem[i:], sep))
		}
	}
	return ""
}
