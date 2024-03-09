package fullpath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFullPath_ToPathLink(t *testing.T) {
	fp := FullPath("/aa/bb/cc")
	assert.Equal(t, fp.ToPathLink(), []PathLink{{"/", "/"}, {"aa/", "/aa/"}, {"bb/", "/aa/bb/"}, {"cc/", "/aa/bb/cc/"}})

	fp = FullPath("/")
	assert.Equal(t, fp.ToPathLink(), []PathLink{{"/", "/"}})
}
