// Package fs provides filesystem-related functions.
package fs

import (
	"errors"
	"os"
	"path/filepath"
)

// Walker provides a convenient interface for iterating over the
// descendents of a filesystem path. It wraps filepath.Walk with
// an interface patterned after bufio.Scanner and sql.Rows.
// Successive calls to the Step method will step through each
// file or directory in the tree, including the root. The files
// are walked in lexical order, which makes the output deterministic
// but means that for very large directories Walker can be inefficient.
// Walker does not follow symbolic links.
//
// Walking continues until Stop is called or there are no more
// filesystem entries.
type Walker struct {
	itemch chan item
	cur    item
	ok     bool // is cur valid?
	retch  chan error
	ret    error
}

type item struct {
	path string
	info os.FileInfo
	err  error
}

// Walk returns a new Walker rooted at root.
func Walk(root string) *Walker {
	w := new(Walker)
	w.itemch = make(chan item)
	w.retch = make(chan error, 1)
	go func() {
		filepath.Walk(root, w.recv)
		close(w.itemch)
	}()
	return w
}

func (w *Walker) recv(path string, info os.FileInfo, err error) error {
	w.itemch <- item{path, info, err}
	return <-w.retch
}

// Step advances the Walker to the next file or directory,
// which will then be available through the Path, Stat,
// and Err methods.
// It returns false when the walk stops, either at the
// end of the tree or from a call to Stop.
func (w *Walker) Step() bool {
	if w.ok {
		w.retch <- w.ret
	}
	w.cur, w.ok = <-w.itemch
	w.ret = nil
	return w.ok
}

// Path returns the path to the most recent file or directory
// visited by a call to Step. It contains the argument to Walk
// as a prefix; that is, if Walk is called with "dir", which is
// a directory containing the file "a", Path will return "dir/a".
func (w *Walker) Path() string {
	return w.cur.path
}

// Stat returns info for the most recent file or directory
// visited by a call to Step.
func (w *Walker) Stat() os.FileInfo {
	return w.cur.info
}

// Err returns the error, if any, for the most recent attempt
// by Step to visit a file or directory. If a directory has
// an error, w will not descend into that directory.
func (w *Walker) Err() error {
	return w.cur.err
}

// SkipDir causes the currently visited directory to be skipped.
// If w is not on a directory, SkipDir has no effect.
func (w *Walker) SkipDir() {
	if w.ok && w.cur.info.IsDir() {
		w.ret = filepath.SkipDir
	}
}

// Stop stops w from visiting any more filesystem entries.
// Subsequent calls to Step will return false.
func (w *Walker) Stop() {
	select {
	case w.retch <- errors.New("stop"):
	default:
	}
}
