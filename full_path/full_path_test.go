package full_path

import (
	"github.com/go-playground/assert/v2"
	"runtime"
	"testing"
)

func TestFullPath_IsLegal(t *testing.T) {
	path := FullPath("/")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/bb/cc")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/bb/./cc")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/bb/../cc")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/bb/.../cc")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("aa/bb/cc/")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/aa/bb/cc/")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/aa/bb//cc")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/./.")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/./aa")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/../bb")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/../../bb")
	assert.Equal(t, path.IsLegal(), false)

	for _, ch := range "<>\\|:*?" {
		path = FullPath("/aa/b" + string(ch) + "b/cc")
		assert.Equal(t, path.IsLegal(), false)
	}
}

func TestFullPath_Clean(t *testing.T) {
	path := FullPath("/aa/bb/./cc")
	if runtime.GOOS == "windows" {
		assert.Equal(t, string(path.Clean()), "\\aa\\bb\\cc")
	} else {
		assert.Equal(t, string(path.Clean()), "/aa/bb/cc")
	}

	path = FullPath("/aa/bb/cc/../dd")
	if runtime.GOOS == "windows" {
		assert.Equal(t, string(path.Clean()), "\\aa\\bb\\dd")
	} else {
		assert.Equal(t, string(path.Clean()), "/aa/bb/dd")
	}

	path = FullPath("/aa/../bb")
	if runtime.GOOS == "windows" {
		assert.Equal(t, string(path.Clean()), "\\bb")
	} else {
		assert.Equal(t, string(path.Clean()), "/bb")
	}

	path = FullPath("/././aa")
	if runtime.GOOS == "windows" {
		assert.Equal(t, string(path.Clean()), "\\aa")
	} else {
		assert.Equal(t, string(path.Clean()), "/aa")
	}

	path = FullPath("/")
	if runtime.GOOS == "windows" {
		assert.Equal(t, string(path.Clean()), "\\")
	} else {
		assert.Equal(t, string(path.Clean()), "/")
	}
}

func TestFullPath_DirAndName(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	dir, name := path.DirAndName()
	assert.Equal(t, dir, "/aa/bb/")
	assert.Equal(t, name, "cc")

	path = FullPath("/aa")
	dir, name = path.DirAndName()
	assert.Equal(t, dir, "/")
	assert.Equal(t, name, "aa")

	path = FullPath("/")
	dir, name = path.DirAndName()
	assert.Equal(t, dir, "/")
	assert.Equal(t, name, "")
}

func TestFullPath_Name(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, path.Name(), "cc")

	path = FullPath("/aa")
	assert.Equal(t, path.Name(), "aa")

	path = FullPath("/")
	assert.Equal(t, path.Name(), "")
}

func TestFullPath_Split(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, path.Split(), []string{"", "aa", "bb", "cc"})

	path = FullPath("/")
	assert.Equal(t, path.Split(), []string{""})
}

func TestFullPath_SplitList(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, path.SplitList(), []string{"/", "/aa", "/aa/bb", "/aa/bb/cc"})

	path = FullPath("/")
	assert.Equal(t, path.SplitList(), []string{"/"})
}
