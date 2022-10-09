package common

import (
	"testing"
)

func Test_Folder(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "/movies/action/matrix.json", out: "/movies/action"},
		{in: "/movies/action/matrix/", out: "/movies/action"},
		{in: "/movies/action/matrix", out: "/movies/action"},
		{in: "/movies", out: ""},
		{in: "movies/action", out: "movies"},
	}

	for _, test := range tests {
		out := Folder(test.in)
		if out != test.out {
			t.Errorf("Expected output to be '%s' but got '%s'", test.out, out)
		}
	}
}
