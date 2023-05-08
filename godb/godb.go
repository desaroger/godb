package godb

import (
	c "godb/common"
	"godb/index"
	"godb/logs"
	s "godb/storage"
	"strings"
)

type Godb struct {
	storage s.Storage
}

func NewGodb(storage s.Storage) *Godb {
	godb := new(Godb)
	godb.storage = storage

	logs.Initialize()

	return godb
}

func (godb *Godb) Set(document c.Document) (c.Document, error) {
	document, err := godb.storage.Set(document)
	if err != nil {
		logs.Error(err, "document.set %v", document)
		return nil, err
	}

	err = index.OnDocumentModified(godb.storage, document)
	if err != nil {
		logs.Error(err, "document.set.updateIndex %v", document)
		return nil, err
	}

	logs.Info("document.set %v", document)

	return document, nil
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

	logs.Info("document.patch %v", document)

	return document, nil
}

func (godb *Godb) List(id string) ([]string, error) {
	ids, err := godb.storage.List(id)
	if err != nil {
		return nil, err
	}

	final_ids := []string{}
	for _, id := range ids {
		if strings.HasSuffix(id, "/") {
			if strings.HasPrefix(id, "_") {
				final_ids = append(final_ids, strings.TrimRight(id, "/"))
			}
			continue
		}

		final_ids = append(final_ids, id)
	}

	return final_ids, nil
}

func (godb *Godb) Delete(id string) error {
	return godb.storage.Delete(id)
}
