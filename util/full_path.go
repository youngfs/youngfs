package util

import (
	"path/filepath"
	"strings"
)

type FullPath string

func NewFullPath(dir, name string) FullPath {
	return FullPath(dir).Child(name)
}

func (fp FullPath) DirAndName() (string, string) {
	dir, fileName := filepath.Split(string(fp))
	fileName = strings.ToValidUTF8(fileName, "?")
	if dir == "/" {
		return dir, fileName
	}
	if len(dir) < 1 {
		return "/", ""
	}
	return dir[:len(dir)-1], fileName
}

func (fp FullPath) FileName() string {
	_, fileName := filepath.Split(string(fp))
	return strings.ToValidUTF8(fileName, "?")
}

func (fp FullPath) Child(name string) FullPath {
	dir := string(fp)
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	if strings.HasSuffix(dir, "/") {
		return FullPath(dir + name)
	}
	return FullPath(dir + "/" + name)
}

func (fp FullPath) Split() []string {
	if fp == "" || fp == "/" {
		return []string{}
	}
	return strings.Split(string(fp)[1:], "/")
}

func Join(names ...string) string {
	return filepath.ToSlash(filepath.Join(names...))
}

func JoinPath(names ...string) FullPath {
	return FullPath(Join(names...))
}
