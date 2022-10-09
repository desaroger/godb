package godb

import (
	"reflect"
	"testing"

	c "godb/common"
	s "godb/storage"
)

func Test_CreateIndex(t *testing.T) {
	storage := &s.FileStorage{Root: t.TempDir()}
	godb := NewGodb(storage)

	index_document := c.NewDocument("movies/_indexes/by_name", "func", "(doc) => ([doc.name, {a: 23, b: doc.id}])")
	err := godb.Set(index_document)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	_, err = godb.Get("movies/_indexes/by_name")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	err = godb.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	err = godb.Set(c.NewDocument("movies/superman", "name", "Superman"))
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	ids, err := godb.List("movies/_indexes/by_name")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if !reflect.DeepEqual(ids, []string{"Matrix", "Superman"}) {
		t.Fatalf("expected ids to be ['Matrix', 'Superman'] but got %s", ids)
	}

	document, err := godb.Get("movies/_indexes/by_name/Matrix")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if !reflect.DeepEqual(document, c.NewDocument("movies/_indexes/by_name/Matrix", "a", float64(23), "b", "movies/matrix")) {
		t.Fatalf("expected indexed document to be {a: 23, b: 'movies/matrix'} but got %s", document)
	}
}
