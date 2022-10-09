package storage

import (
	"errors"
	"reflect"
	"testing"

	c "godb/common"
)

func Test_canCRUDDocuments(t *testing.T) {
	storage := &FileStorage{Root: t.TempDir()}

	// Create matrix movie
	document := c.NewDocument("movies/matrix", "name", "Matrix")
	err := storage.Set(document)
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	// Check if matrix movie exists
	exists1, err := storage.Exists("movies/matrix")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if !exists1 {
		t.Fatalf("expected file exists")
	}

	// Check that an unknown movie does not exist
	exists2, err := storage.Exists("movies/my_girl")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if exists2 {
		t.Fatalf("expected document to not exist but exists")
	}

	// Check matrix movie has the correct data
	document_got, err := storage.Get("movies/matrix")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if document_got == nil {
		t.Fatalf("expected document but got nil")
	}
	if document_got["id"] != "movies/matrix" {
		t.Fatalf("expected document.id='movies/matrix' but got '%s'", document_got["id"])
	}
	if document_got["name"] != "Matrix" {
		t.Fatalf("expected document.name='Matrix' but got '%s'", document_got["name"])
	}

	// Add a description to matrix movie
	err = storage.Patch(c.NewDocument("movies/matrix", "desc", "it's about a guy..."))
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}

	// Check the description is written
	document_got2, err := storage.Get("movies/matrix")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if document_got2 == nil {
		t.Fatalf("expected document but got nil")
	}
	// if document_got2["id"] != "movies/matrix" {
	// 	t.Fatalf("expected document.id='movies/matrix' but got '%s'", document_got2["id"])
	// }
	if document_got2["name"] != "Matrix" {
		t.Fatalf("expected document.name='Matrix' but got '%s'", document_got2["name"])
	}
	if document_got2["desc"] != "it's about a guy..." {
		t.Fatalf("expected document.desc='it's about a guy...' but got '%s'", document_got2["desc"])
	}

	// List root and movies
	ids, err := storage.List("movies")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if !reflect.DeepEqual(ids, []string{"matrix"}) {
		t.Fatalf("expected ids to be [action] but got %s", ids)
	}
	ids2, err := storage.List("unexisting")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if ids == nil {
		t.Fatalf("expected empty list, but got nil")
	}
	if len(ids2) != 0 {
		t.Fatalf("expected empty list, but got '%s'", ids2)
	}

	// Check that does not return nil
	storage.Set(c.NewDocument("books/romance/pride_and_prejudice", "name", "Pride and prejudice"))
	ids3, err := storage.List("books")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	if ids3 == nil {
		t.Fatalf("expected empty list, but got nil")
	}
	if !reflect.DeepEqual(ids3, []string{"romance/"}) {
		t.Fatalf("expected ids to be ['romance/'], but got '%s'", ids3)
	}

	// Delete file
	err = storage.Delete("movies/matrix")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	err = storage.Delete("movies/matrix")
	if !errors.Is(err, c.ErrDocumentDoestNotExist) {
		t.Fatalf("expecting error ErrDocumentDoestNotExist but got %s", err)
	}
	ids4, err := storage.List("movies")
	if err != nil {
		t.Fatalf("unexpected error '%s'", err)
	}
	c.D("ids4", ids4)
	if !reflect.DeepEqual(ids4, []string{}) {
		t.Fatalf("expected ids to be [] but got %s", ids4)
	}
}
