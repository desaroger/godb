package common

import (
	"fmt"
	"path/filepath"

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
	folder, _ := filepath.Split(path)

	return folder
}
