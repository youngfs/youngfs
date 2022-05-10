package set

import (
	"unicode/utf8"
)

type Set string

func (set Set) IsLegal() bool {
	if len(set) == 0 {
		return false
	}

	if !utf8.ValidString(string(set)) {
		return false
	}

	for _, u := range set {
		if u == '/' {
			return false
		}
	}

	return true
}
