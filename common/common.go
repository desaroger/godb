package common

import (
	"encoding/json"
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

func DeepClone[T any](x T) T {
	bytes, err := json.Marshal(x)
	if err != nil {
		panic("XXX")
	}

	clone := new(T)
	// clone := reflect.New(reflect.TypeOf(x))
	json.Unmarshal(bytes, clone)

	return *clone
}
