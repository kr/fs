// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"path/filepath"
	"runtime"
	"testing"
)

type JoinTest struct {
	elem []string
	path string
}

var jointests = []JoinTest{
	// zero parameters
	{[]string{}, ""},

	// one parameter
	{[]string{""}, ""},
	{[]string{"a"}, "a"},
	// one parameter
	{[]string{""}, ""},
	{[]string{"a"}, "a"},

	// two parameters
	{[]string{"a", "b"}, "a/b"},
	{[]string{"a", ""}, "a"},
	{[]string{"", "b"}, "b"},
	{[]string{"/", "a"}, "/a"},
	{[]string{"/", ""}, "/"},
	{[]string{"a/", "b"}, "a/b"},
	{[]string{"a/", ""}, "a"},
	{[]string{"", ""}, ""},
}

var winjointests = []JoinTest{
	{[]string{`directory`, `file`}, `directory\file`},
	{[]string{`C:\Windows\`, `System32`}, `C:\Windows\System32`},
	{[]string{`C:\Windows\`, ``}, `C:\Windows`},
	{[]string{`C:\`, `Windows`}, `C:\Windows`},
	{[]string{`C:`, `Windows`}, `C:\Windows`},
	{[]string{`\\host\share`, `foo`}, `\\host\share\foo`},
	{[]string{`//host/share`, `foo/bar`}, `\\host\share\foo\bar`},
}

// join takes a []string and passes it to Join.
func join(elem []string, args ...string) string {
	args = elem
	return Join(new(fs), args...)
}

func TestJoin(t *testing.T) {
	if runtime.GOOS == "windows" {
		jointests = append(jointests, winjointests...)
	}
	for _, test := range jointests {
		if p := join(test.elem); p != filepath.FromSlash(test.path) {
			t.Errorf("join(%q) = %q, want %q", test.elem, p, test.path)
		}
	}
}
