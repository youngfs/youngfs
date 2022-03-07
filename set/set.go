package set

import (
	"unicode/utf8"
)

type Set string

func (set Set) IsLegal() bool {
	if !utf8.ValidString(string(set)) {
		return false
	}

	return true
}
