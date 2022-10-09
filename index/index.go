package index

import (
	c "godb/common"
	s "godb/storage"
)

type Index struct {
	Id   string `json:"id"`
	Func string `json:"func"`
}

func indexFromDocument(document c.Document) Index {
	return Index{
		Id:   document["id"].(string),
		Func: document["func"].(string),
	}
}

func load_indexes(storage s.Storage, folder string) ([]Index, error) {
	indexes_folder := c.J(folder, "_indexes")
	indexes_simple_ids, err := storage.List(indexes_folder)
	if err != nil {
		return nil, err
	}

	var indexes []Index
	for _, index_simple_id := range indexes_simple_ids {
		index_id := c.J(indexes_folder, index_simple_id)
		index_document, err := storage.Get(index_id)
		if err != nil {
			return nil, err
		}

		index := indexFromDocument(index_document)
		indexes = append(indexes, index)
	}

	return indexes, nil
}

func OnDocumentModified(storage s.Storage, document c.Document) error {
	document_id, err := document.GetId()
	if err != nil {
		return err
	}

	document_folder := c.Folder(document_id)
	indexes, err := load_indexes(storage, document_folder)
	if err != nil {
		return err
	}

	document_simple_ids, err := storage.List(document_folder)
	if err != nil {
		return err
	}

	for _, index := range indexes {
		err = storage.DeleteFolder(index.Id)
		if err != nil {
			return err
		}
	}

	for _, document_simple_id := range document_simple_ids {
		document_id := c.J(document_folder, document_simple_id)
		document, err := storage.Get(document_id)
		if err != nil {
			// TODO ya veremos que hacemos
			return err
		}

		for _, index := range indexes {
			evaluation_id, evaluation_content, err := evaluate(document, index.Func)
			if err != nil {
				// TODO ya veremos que hacemos
				return err
			}
			if evaluation_id == "" {
				continue
			}

			evaluation_document_id := c.J(index.Id, evaluation_id)
			evaluation_content["id"] = evaluation_document_id

			err = storage.Set(evaluation_content)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
