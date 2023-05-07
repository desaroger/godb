package storage

import (
	"errors"
	"reflect"
	"testing"

	c "godb/common"
)

func getStorages(tempDir string) map[string]Storage {
	return map[string]Storage{
		"FileStorage":   &FileStorage{Root: tempDir},
		"MemoryStorage": &MemoryStorage{},
	}
}

func testEachStorage(t *testing.T, f func(*testing.T, Storage)) {
	for storageName, storage := range getStorages(t.TempDir()) {
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
		if !objectsDeepEqual(document, expected) {
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
		if !objectsDeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{"matrix"}
		if !objectsDeepEqual(ids, expected) {
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
		if !objectsDeepEqual(document, expected) {
			t.Fatalf("expected 'movies/matrix' document to equal %v but got %v", expected, document)
		}

		// Patch the description and add a year
		_, err = storage.Patch(c.NewDocument("movies/matrix", "desc", "It's about a guy...", "year", 1999))
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}

		// First check if the file matches what we wrote
		document, err = storage.Get("movies/matrix")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = c.NewDocument("movies/matrix", "name", "Matrix", "desc", "It's about a guy...", "year", 1999)
		if !objectsDeepEqual(document, expected) {
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
		if !objectsDeepEqual(ids, expected) {
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
		if !objectsDeepEqual(ids, expected) {
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
		if !objectsDeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// List movies
		ids, err = storage.List("movies")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{"matrix"}
		if !objectsDeepEqual(ids, expected) {
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
		if !objectsDeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}

		// List root
		ids, err = storage.List("")
		if err != nil {
			t.Fatalf("unexpected error '%s'", err)
		}
		expected = []string{}
		if !objectsDeepEqual(ids, expected) {
			t.Fatalf("expected list ids to equal %v but got %v", expected, ids)
		}
	})
}

func BenchmarkStorages(b *testing.B) {
	for storageName, storage := range getStorages(b.TempDir()) {
		b.Run(storageName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				storage.Set(c.NewDocument("movies/matrix", "name", "Matrix"))
				storage.List("movies")
				storage.Get("movies/matrix")
				storage.Delete("movies/matrix")
			}
		})
	}
}

func objectsDeepEqual(x any, y any) bool {
	xx := c.DeepClone(x)
	yy := c.DeepClone(y)

	return reflect.DeepEqual(xx, yy)
}
