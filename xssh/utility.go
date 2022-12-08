package xssh

import (
	"os"
	"strings"
)

type FileInfo struct {
	Name string
	Path string
}

func (f *FileInfo) IsDir() bool {
	return strings.HasSuffix(f.Name, "/") || (f.Name == "" && strings.HasSuffix(f.Path, "/"))
}

func Dir(path string) string {
	split := strings.Split(path, "/")
	if len(split) > 1 {
		return strings.Join(split[:len(split)-1], "/")
	} else {
		return "/"
	}
}

func FileName(path string) string {
	split := strings.Split(path, "/")
	if len(split) > 1 {
		return split[len(split)-1]
	} else {
		return path
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func Command(name string, arg ...string) string {
	if arg != nil {
		name += " " + strings.Join(arg, " ")
	}
	return name
}
