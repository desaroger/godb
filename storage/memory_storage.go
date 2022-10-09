package storage

import (
	"errors"
	"strings"

	c "godb/common"
)

type MemoryStorage struct {
	Storage
	data map[string]c.Document
}

func (ms *MemoryStorage) Get(id string) (c.Document, error) {
	ms.ensureData()
	document := ms.data[id]
	if document == nil {
		return nil, c.ErrDocumentDoestNotExist
	}

	return document, nil
}

func (ms *MemoryStorage) Set(document c.Document) error {
	ms.ensureData()
	document_id, err := document.GetId()
	if err != nil {
		return err
	}
	ms.data[document_id] = document

	return nil
}

func (ms *MemoryStorage) Patch(document c.Document) error {
	ms.ensureData()
	document_id, err := document.GetId()
	if err != nil {
		return err
	}

	existing_document, err := ms.Get(document_id)
	if err != nil {
		if !errors.Is(err, c.ErrDocumentDoestNotExist) {
			return err
		}
	}
	if existing_document == nil {
		existing_document = c.NewDocument(document_id)
	}

	existing_document = c.DeepClone(existing_document)
	existing_document.Patch(document)

	return ms.Set(existing_document)
}

func (ms *MemoryStorage) Exists(id string) (bool, error) {
	ms.ensureData()
	_, exists := ms.data[id]

	return exists, nil
}

func (ms *MemoryStorage) List(folder string) ([]string, error) {
	ms.ensureData()

	folder = strings.Trim(folder, "/")
	simple_ids := []string{}
	for document_id := range ms.data {
		if folder != "" && !strings.HasPrefix(document_id, folder) {
			continue
		}
		relative_id := strings.TrimPrefix(document_id, folder)
		relative_id = strings.Trim(relative_id, "/")
		relative_id_parts := strings.Split(relative_id, "/")

		if len(relative_id_parts) == 1 {
			simple_ids = append(simple_ids, relative_id_parts[0])
		} else {
			simple_ids = append(simple_ids, relative_id_parts[0]+"/")
		}
	}

	return simple_ids, nil
}

func (ms *MemoryStorage) Delete(id string) error {
	ms.ensureData()

	delete(ms.data, id)

	return nil
}

func (ms *MemoryStorage) ensureData() {
	if ms.data == nil {
		ms.data = map[string]c.Document{}
	}
}
