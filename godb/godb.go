package godb

import (
	c "godb/common"
	"godb/index"
	"godb/logs"
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
		logs.Error(err, "document.set %w", document)
		return err
	}

	err = index.OnDocumentModified(godb.storage, document)
	if err != nil {
		logs.Error(err, "document.set.updateIndex %w", document)
		return err
	}

	logs.Info("document.set %s", document.GetIdOrNil())

	return nil
}

func (godb *Godb) Get(id string) (c.Document, error) {
	return godb.storage.Get(id)
}

func (godb *Godb) Patch(document c.Document) (c.Document, error) {
	document, err := godb.storage.Patch(document)
	if err != nil {
		logs.Error(err, "document.patch %v", document)
		return nil, err
	}

	err = index.OnDocumentModified(godb.storage, document)
	if err != nil {
		logs.Error(err, "document.patch.updateIndex %v", document)
		return nil, err
	}

	logs.Info("document.patch %s", document.GetIdOrNil())

	return document, nil
}

func (godb *Godb) List(id string) ([]string, error) {
	return godb.storage.List(id)
}

func (godb *Godb) Delete(id string) error {
	return godb.storage.Delete(id)
}
