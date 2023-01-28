package full_path

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestFullPath_ToPathLink(t *testing.T) {
	fp := FullPath("/aa/bb/cc")
	assert.Equal(t, fp.ToPathLink(), []PathLink{{"/", "/"}, {"aa/", "/aa/"}, {"bb/", "/aa/bb/"}, {"cc/", "/aa/bb/cc/"}})

	fp = FullPath("/")
	assert.Equal(t, fp.ToPathLink(), []PathLink{{"/", "/"}})
}
