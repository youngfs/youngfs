package util

import "strings"

// PrefixEnd returns the immediate next byte slice lexicographically
// after the given prefix. The result can be used as the exclusive end
// key when performing a range scan with the provided prefix.
//
// For example:
//   - For prefix "abc", the result might be "abd".
//   - For prefix "ab\xff", the result might be "ac\x00".
//   - For a prefix that ends in multiple 0xff bytes like "ab\xff\xff",
//     the result would be "ac\x00\x00".
//
// If the given prefix consists entirely of 0xff bytes, the function
// returns nil, indicating there's no valid end key for such prefix.
//
// It is not guaranteed that the returned byte slice is a valid key in
// the user's data domain, but it is suitable for range scan operations.
//
// Params:
//   - prefix: The byte slice for which the end key is desired.
//
// Returns:
//   - A byte slice representing the end key, or nil if the prefix is
//     composed entirely of 0xff bytes.
func PrefixEnd(prefix []byte) []byte {
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		if end[i] != 0xff {
			end[i]++
			return end
		}
		end[i] = 0
	}
	return nil
}

func TrimEtag(etag string) string {
	if strings.HasPrefix(etag, "\"") {
		etag = etag[1:]
	}
	if strings.HasSuffix(etag, "\"") {
		etag = etag[:len(etag)-1]
	}
	return etag
}

func QuotEtag(etag string) string {
	if !strings.HasPrefix(etag, "\"") {
		etag = "\"" + etag
	}
	if !strings.HasSuffix(etag, "\"") {
		etag = etag + "\""
	}
	return etag
}
