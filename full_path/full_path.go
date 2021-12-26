package full_path

import (
	"path/filepath"
	"strings"
)

type FullPath string

/*
	example:	/aa/bb/cc
	dir:		/aa/bb/
	name:		cc

	example:	/aa
	dir:		/
	name:		aa

	example:	/aa/bb/cc/../dd
	dir:		/aa/bb/
	name:		dd

	example:	/
	dir:		/
	name:

	use IsLegal before other functions
*/

func (fp FullPath) IsLegal() bool {
	if fp == "" {
		return false
	}
	if fp == "/" {
		return true
	}
	if fp[0] != '/' || fp[len(fp)-1] == '/' {
		return false
	}
	set := make(map[int32]bool)
	for _, ch := range "<>\\|:*?" {
		set[ch] = true
	}
	for _, ch := range fp {
		if set[ch] == true {
			return false
		}
	}
	if strings.Contains(string(fp), "//") || strings.Contains(string(fp), "...") {
		return false
	}
	cnt := 0
	for _, str := range strings.Split(string(fp), "/") {
		if str == ".." {
			cnt--
		} else if str != "." && str != "" {
			cnt++
		}
		if cnt < 0 {
			return false
		}
	}
	if cnt <= 0 {
		return false
	}
	return true
}

func (fp FullPath) Clean() FullPath {
	// os:  Windows: \ Linux: /
	return FullPath(filepath.Clean(string(fp)))
}

func (fp FullPath) DirAndName() (string, string) {
	// fp = fp.Clean()
	dir, name := filepath.Split(string(fp))
	name = strings.ToValidUTF8(name, "?")
	return dir, name
}

func (fp FullPath) Name() string {
	// fp = fp.Clean()
	_, name := filepath.Split(string(fp))
	return strings.ToValidUTF8(name, "?")
}

func (fp FullPath) Split() []string {
	if fp == "/" {
		return []string{""}
	}
	return strings.Split(string(fp), "/")
}

func (fp FullPath) SplitList() []string {
	list := fp.Split()
	for i := 1; i < len(list); i++ {
		list[i] = list[i-1] + "/" + list[i]
	}
	list[0] = "/"
	return list
}
