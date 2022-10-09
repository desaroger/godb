package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func S(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func D(a ...interface{}) {
	spew.Dump(a...)
}

func J(elem ...string) string {
	return filepath.Join(elem...)
}

func Folder(path string) string {
	path = strings.TrimRight(path, "/")
	folder, _ := filepath.Split(path)
	folder = strings.TrimRight(folder, "/")

	return folder
}
