package full_path

import (
	"github.com/go-playground/assert/v2"
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

	path = FullPath("aa/bb/cc")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("aa/bb/cc/")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/aa/bb/cc/")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/aa/bb//cc")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/.")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/./.")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/./aa")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/../bb")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/..")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa/../../bb")
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/aa/bb/../..")
	assert.Equal(t, path.IsLegal(), true)

	for _, ch := range "<>\\|:*?" {
		path = FullPath("/aa/b" + string(ch) + "b/cc")
		assert.Equal(t, path.IsLegal(), false)
	}

	path = FullPath("/测试")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/ÄäÖöÜüẞß")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/あいうえお")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/aa bb")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/" + string([]byte{0xff, 0xfe, 0xfd}))
	assert.Equal(t, path.IsLegal(), false)

	path = FullPath("/  aa/\tbb/  cc")
	assert.Equal(t, path.IsLegal(), true)

	path = FullPath("/  aa dd/\tbb/  cc")
	assert.Equal(t, path.IsLegal(), true)
}

func TestFullPath_IsLegalObjectName(t *testing.T) {
	path := FullPath("/")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/aa/bb/cc")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa/bb/./cc")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa/bb/../cc")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa/bb/.../cc")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("aa/bb/cc")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("aa/bb/cc/")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/aa/bb/cc/")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/aa/bb//cc")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/.")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/./.")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/./aa")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa/../bb")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa/..")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/aa/../../bb")
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/aa/bb/../..")
	assert.Equal(t, path.IsLegalObjectName(), false)

	for _, ch := range "<>\\|:*?" {
		path = FullPath("/aa/b" + string(ch) + "b/cc")
		assert.Equal(t, path.IsLegalObjectName(), false)
	}

	path = FullPath("/测试")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/ÄäÖöÜüẞß")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/あいうえお")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/aa bb")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/" + string([]byte{0xff, 0xfe, 0xfd}))
	assert.Equal(t, path.IsLegalObjectName(), false)

	path = FullPath("/  aa/\tbb/  cc")
	assert.Equal(t, path.IsLegalObjectName(), true)

	path = FullPath("/  aa dd/\tbb/  cc")
	assert.Equal(t, path.IsLegalObjectName(), true)
}

func TestFullPath_Clean(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, string(path.Clean()), "/aa/bb/cc")

	path = FullPath("/aa/bb/./cc")
	assert.Equal(t, string(path.Clean()), "/aa/bb/cc")

	path = FullPath("/aa/bb/cc/../dd")
	assert.Equal(t, string(path.Clean()), "/aa/bb/dd")

	path = FullPath("/aa/../bb")
	assert.Equal(t, string(path.Clean()), "/bb")

	path = FullPath("/aa/bb/../..")
	assert.Equal(t, string(path.Clean()), "/")

	path = FullPath("/././aa")
	assert.Equal(t, string(path.Clean()), "/aa")

	path = FullPath("/  aa/\tbb/  cc")
	assert.Equal(t, string(path.Clean()), "/aa/bb/cc")

	path = FullPath("/  aa dd/\tbb/  cc")
	assert.Equal(t, string(path.Clean()), "/aa dd/bb/cc")

	// illegal full path object name
	path = FullPath("/")
	assert.Equal(t, string(path.Clean()), "/")

	path = FullPath("/.")
	assert.Equal(t, string(path.Clean()), "/")

	path = FullPath("/./.")
	assert.Equal(t, string(path.Clean()), "/")

	path = FullPath("/aa/..")
	assert.Equal(t, string(path.Clean()), "/")
}

func TestFullPath_DirAndName(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	dir, name := path.DirAndName()
	assert.Equal(t, dir, FullPath("/aa/bb"))
	assert.Equal(t, name, "cc")

	path = FullPath("/aa")
	dir, name = path.DirAndName()
	assert.Equal(t, dir, FullPath("/"))
	assert.Equal(t, name, "aa")

	path = FullPath("/")
	dir, name = path.DirAndName()
	assert.Equal(t, dir, FullPath("/"))
	assert.Equal(t, name, "")
}

func TestFullPath_Dir(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, path.Dir(), FullPath("/aa/bb"))

	path = FullPath("/aa")
	assert.Equal(t, path.Dir(), FullPath("/"))

	path = FullPath("/")
	assert.Equal(t, path.Dir(), FullPath("/"))
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
	assert.Equal(t, path.Split(), []FullPath{"", "aa", "bb", "cc"})

	path = FullPath("/")
	assert.Equal(t, path.Split(), []FullPath{""})
}

func TestFullPath_SplitList(t *testing.T) {
	path := FullPath("/aa/bb/cc")
	assert.Equal(t, path.SplitList(), []FullPath{"/", "/aa", "/aa/bb", "/aa/bb/cc"})

	path = FullPath("/")
	assert.Equal(t, path.SplitList(), []FullPath{"/"})
}
