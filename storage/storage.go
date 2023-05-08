package storage

import (
	c "godb/common"
)

type Storage interface {
	Get(id string) (c.Document, error)
	Set(document c.Document) (c.Document, error)
	Patch(document c.Document) (c.Document, error)
	Exists(id string) (bool, error)
	List(folder string) ([]string, error)
	Delete(id string) error
	DeleteFolder(folder string) error
}
