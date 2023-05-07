package godb

import (
	c "godb/common"
	"godb/index"
	s "godb/storage"
)

type Godb struct {
	storage s.Storage
}

func NewGodb(storage s.Storage) *Godb {
	godb := new(Godb)
	godb.storage = storage

	return godb
}

func (godb *Godb) Set(document c.Document) error {
	err := godb.storage.Set(document)
	if err != nil {
		return err
	}

	return index.OnDocumentModified(godb.storage, document)
}

func (godb *Godb) Get(id string) (c.Document, error) {
	return godb.storage.Get(id)
}

func (godb *Godb) Patch(document c.Document) (c.Document, error) {
	document, err := godb.storage.Patch(document)
	if err != nil {
		return nil, err
	}

	err = index.OnDocumentModified(godb.storage, document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (godb *Godb) List(id string) ([]string, error) {
	return godb.storage.List(id)
}

func (godb *Godb) Delete(id string) error {
	return godb.storage.Delete(id)
}
