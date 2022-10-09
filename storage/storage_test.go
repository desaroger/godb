package storage

import (
	"errors"
	"reflect"
	"testing"

	c "godb/common"
)

func testEachStorage(t *testing.T, f func(*testing.T, Storage)) {
	storages := map[string]Storage{
		"FileStorage": &FileStorage{Root: t.TempDir()},
	}

	for storageName, storage := range storages {
		t.Run(storageName, func(tt *testing.T) {
			f(tt, storage)
		})
	}
}

func Test_CanCreateDocuments(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		exists, err := storage.Exists("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		if exists != true {
			t.Fatalf("expected 'movies/matrix' to exists but didn't")
		}

		exists, err = storage.Exists("movies/nope")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		if exists != false {
			t.Fatalf("expected 'movies/nope' to not exist but exists")
		}
	})
}

func Test_CanGetDocuments(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		document, err := storage.Get("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected := c.NewDocument("movies/matrix", "name", "Matrix")
		if !reflect.DeepEqual(document, expected) {
			t.Fatalf("expected 'movies/matrix' document to equal %v but got %v", expected, document)
		}

		_, err = storage.Get("movies/nope")
		if !errors.Is(err, c.ErrDocumentDoestNotExist) {
			t.Fatalf("unexpected error 'DocumentDoesNotExist' but got '%s'", err)
		}
	})
}

func Test_CanListDocuments(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		ids, err := storage.List("")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		if ids == nil {
			t.Fatalf("expected list ids to be an array but got nil")
		}
		expected := []string{"movies/"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{"matrix"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}
	})
}

func Test_CanPatchDocuments(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix", "desc", "wrong"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// First check if the file matches what we wrote
		document, err := storage.Get("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected := c.NewDocument("movies/matrix", "name", "Matrix", "desc", "wrong")
		if !reflect.DeepEqual(document, expected) {
			t.Fatalf("expected 'movies/matrix' document to equal %v but got %v", expected, document)
		}

		// Patch the description and add a year
		err = storage.Patch(c.NewDocument("movies/matrix", "desc", "It's about a guy...", "year", 1999))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// First check if the file matches what we wrote
		document, err = storage.Get("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = c.NewDocument("movies/matrix", "name", "Matrix", "desc", "It's about a guy...", "year", float64(1999))
		if !reflect.DeepEqual(document, expected) {
			t.Fatalf("expected 'movies/matrix' document to equal %v but got %v", expected, document)
		}
	})
}

func Test_CanDeleteDocuments(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		err = storage.Set(c.NewDocument("movies/pride_and_prejudice", "name", "Pride and prejudice"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// List movies
		ids, err := storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected := []string{"matrix", "pride_and_prejudice"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// Delete matrix movie
		err = storage.Delete("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// List movies to see if has been removed
		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{"pride_and_prejudice"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}
	})
}

func Test_DeleteRemovesEmptyFolders(t *testing.T) {
	testEachStorage(t, func(t *testing.T, storage Storage) {
		err := storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// List root
		ids, err := storage.List("")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected := []string{"movies/"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// List movies
		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{"matrix"}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// Delete matrix movie
		err = storage.Delete("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// List movies to see if has been removed
		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// List root
		ids, err = storage.List("")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{}
		if !reflect.DeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}
	})
}
