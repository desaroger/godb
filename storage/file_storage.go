package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	c "godb/common"
)

type FileStorage struct {
	Storage
	Root string
}

func (fs *FileStorage) Get(id string) (c.Document, error) {
	path := fs.resolvePath(id + ".json")

	return fs.fileGet(path)
}

func (fs *FileStorage) Set(document c.Document) (c.Document, error) {
	document_id, err := document.GetId()
	if err != nil {
		return nil, err
	}
	if document_id == "" {
		return nil, c.ErrInvalidId
	}

	path := fs.resolvePath(document_id + ".json")

	return fs.fileSet(path, document)
}

func (fs *FileStorage) Patch(document c.Document) (c.Document, error) {
	document_id, err := document.GetId()
	if err != nil {
		return nil, err
	}
	path := fs.resolvePath(document_id + ".json")

	existing_document, err := fs.fileGet(path)
	if err != nil {
		if !errors.Is(err, c.ErrDocumentDoestNotExist) {
			return nil, err
		}
	}
	if existing_document == nil {
		existing_document = c.NewDocument(document_id)
	}

	existing_document.Patch(document)

	return fs.fileSet(path, existing_document)
}

func (fs *FileStorage) Exists(id string) (bool, error) {
	path := fs.resolvePath(id + ".json")

	return fs.fileExists(path)
}

func (fs *FileStorage) List(folder string) ([]string, error) {
	path := fs.resolvePath(folder)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	simple_ids := []string{}
	for _, f := range files {
		simple_id := f.Name()
		if strings.HasSuffix(f.Name(), ".json") {
			simple_id = strings.TrimSuffix(simple_id, ".json")
		} else {
			simple_id += "/"
		}
		simple_ids = append(simple_ids, simple_id)
	}

	return simple_ids, nil
}

func (fs *FileStorage) Delete(id string) error {
	path := fs.resolvePath(id + ".json")

	err := fs.ensureFolder(path)
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return c.ErrDocumentDoestNotExist
		}
		return err
	}

	return fs.removeEmptyFolders(id)
}

func (fs *FileStorage) DeleteFolder(folder string) error {
	path := fs.resolvePath(folder)

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return fs.removeEmptyFolders(folder)
}

func (fs *FileStorage) fileGet(path string) (c.Document, error) {
	err := fs.ensureFolder(path)
	if err != nil {
		return nil, err
	}

	document_bytes, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, c.ErrDocumentDoestNotExist
		}
		return nil, err
	}

	document := c.Document{}
	err = json.Unmarshal(document_bytes, &document)
	if err != nil {
		return nil, err
	}

	return document, err
}

func (fs *FileStorage) fileCreate(path string, document c.Document) error {
	err := fs.ensureFolder(path)

	// open file for creation
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return c.ErrDocumentAlreadyExists
		}
		return err
	}
	defer file.Close()

	// serialize
	document_bytes, err := json.Marshal(document)
	if err != nil {
		return err
	}

	// write
	_, err = file.Write(document_bytes)

	return err
}

func (fs *FileStorage) fileSet(path string, document c.Document) (c.Document, error) {
	backup_path := path + ".backup"

	err := fs.ensureFolder(path)
	if err != nil {
		return nil, err
	}

	err = fs.fileCreate(backup_path, document)
	if err != nil {
		return nil, err
	}

	err = os.Remove(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	err = os.Rename(backup_path, path)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (fs *FileStorage) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (fs *FileStorage) resolvePath(relative_path string) string {
	if strings.Contains(relative_path, "..") {
		panic("NO!! '" + relative_path + "'")
	}
	return filepath.Join(fs.Root, relative_path)
}

func (fs *FileStorage) ensureFolder(path string) error {
	folder := c.Folder(path)

	return os.MkdirAll(folder, os.ModePerm)
}

func (fs *FileStorage) removeEmptyFolders(id string) error {
	if id == "" {
		return nil
	}

	path := fs.resolvePath(id)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			parentFolder := c.Folder(id)
			return fs.removeEmptyFolders(parentFolder)
		}
		return err
	}
	if len(files) != 0 {
		return nil
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	parentFolder := c.Folder(id)
	return fs.removeEmptyFolders(parentFolder)
}
