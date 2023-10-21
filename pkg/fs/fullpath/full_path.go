package fullpath

import (
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type FullPath string

/*
	example:	/aa/bb/cc
	dir:		/aa/bb
	name:		cc

	example:	/aa
	dir:		/
	name:		aa

	example:	/aa/bb/cc/../dd
	dir:		/aa/bb
	name:		dd

	example:	/
	dir:		/
	name:

	use IsLegal before other functions
*/

func (fp FullPath) IsLegal() bool {
	if !utf8.ValidString(string(fp)) {
		return false
	}
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
	return true
}

func (fp FullPath) IsLegalObjectName() bool {
	if !utf8.ValidString(string(fp)) {
		return false
	}

	if fp == "" {
		return false
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
	if fp == "/" {
		return fp
	}
	list := strings.Split(string(fp), "/")
	list = list[1:]
	retList := make([]string, 0)
	for _, dir := range list {
		if dir == ".." {
			retList = retList[:len(retList)-1]
		} else if dir != "." {
			retList = append(retList, dir)
		}
	}

	ret := ""
	for _, dir := range retList {
		p := 0
		for i, u := range dir {
			if u != ' ' && u != '\t' {
				p = i
				break
			}
		}
		dir = dir[p:]
		ret += "/" + dir
	}
	if ret == "" {
		ret = "/"
	}

	return FullPath(ret)
}

func (fp FullPath) DirAndName() (FullPath, string) {
	// fp = fp.Clean()
	dir, name := filepath.Split(string(fp))
	if dir != "/" {
		dir = dir[:len(dir)-1]
	}
	return FullPath(dir), name
}

func (fp FullPath) Dir() FullPath {
	// fp = fp.Clean()
	dir, _ := filepath.Split(string(fp))
	if dir != "/" {
		dir = dir[:len(dir)-1]
	}
	return FullPath(dir)
}

func (fp FullPath) Name() string {
	// fp = fp.Clean()
	_, name := filepath.Split(string(fp))
	return strings.ToValidUTF8(name, "?")
}

func (fp FullPath) Split() []FullPath {
	if fp == "/" {
		return []FullPath{""}
	}

	list := strings.Split(string(fp), "/")

	ret := make([]FullPath, len(list))
	for i, v := range list {
		ret[i] = FullPath(v)
	}
	return ret
}

func (fp FullPath) SplitList() []FullPath {
	list := fp.Split()
	for i := 1; i < len(list); i++ {
		list[i] = list[i-1] + "/" + list[i]
	}
	list[0] = "/"
	return list
}
