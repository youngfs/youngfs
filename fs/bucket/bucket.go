package bucket

import (
	"unicode/utf8"
)

type Bucket string

func (bkt Bucket) IsLegal() bool {
	if len(bkt) == 0 {
		return false
	}

	if !utf8.ValidString(string(bkt)) {
		return false
	}

	for _, u := range bkt {
		if u == '/' {
			return false
		}
	}

	return true
}
