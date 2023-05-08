package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	c "godb/common"
	"godb/godb"
	s "godb/storage"
)

type httpJsonApi struct {
	godb *godb.Godb
}

func NewHttpJsonApi(rootFolder string) *httpJsonApi {
	storage := &s.FileStorage{Root: rootFolder}
	return &httpJsonApi{
		godb: godb.NewGodb(storage),
	}
}

func (api *httpJsonApi) Start(addr string) {
	http.HandleFunc("/", api.handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (api *httpJsonApi) handle(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	query := r.URL.Query()

	response := api.handleInner(path, query)

	if err, ok := response.(error); ok {
		response = c.Document{
			"error": err.Error(),
		}
	}

	json.NewEncoder(w).Encode(response)
}

func (api *httpJsonApi) handleInner(path string, query url.Values) any {
	for key := range query {
		if key == "id" {
			return errors.New("id needs to be set in the url")
		}
	}

	itemsList := strings.Split(path, "/")
	lastItem := itemsList[len(itemsList)-1]

	switch lastItem {
	case "_set":
		id := c.J(itemsList[:len(itemsList)-1]...)
		document := queryToDocument(id, query)
		return api.set(document)
	case "_patch":
		id := c.J(itemsList[:len(itemsList)-1]...)
		fmt.Println("id " + id)

		document := queryToDocument(id, query)
		return api.patch(document)
	case "_delete":
		id := c.J(itemsList[:len(itemsList)-1]...)
		return api.delete(id)
	case "_list":
		id := c.J(itemsList[:len(itemsList)-1]...)
		return api.list(id)
	}

	return api.get(path)
}

func (api *httpJsonApi) get(id string) any {
	if id == "" {
		id = "/"
	}

	document, err := api.godb.Get(id)
	if err != nil {
		if !errors.Is(err, c.ErrDocumentDoestNotExist) {
			return err
		}
	}

	ids, err := api.godb.List(id)
	if err != nil {
		return err
	}

	return c.Document{
		"document": document,
		"childIds": ids,
	}
}

func (api *httpJsonApi) set(document c.Document) any {
	document, err := api.godb.Set(document)
	if err != nil {
		return err
	}

	ids, err := api.godb.List(document.GetIdOrNil())
	if err != nil {
		return err
	}

	return c.Document{
		"document": document,
		"childIds": ids,
	}
}

func (api *httpJsonApi) patch(document c.Document) any {
	document, err := api.godb.Patch(document)
	if err != nil {
		return err
	}

	ids, err := api.godb.List(document.GetIdOrNil())
	if err != nil {
		return err
	}

	return c.Document{
		"document": document,
		"childIds": ids,
	}
}

func (api *httpJsonApi) list(id string) any {
	ids, err := api.godb.List(id)
	if err != nil {
		return err
	}

	return ids
}

func (api *httpJsonApi) delete(id string) any {
	err := api.godb.Delete(id)
	if err != nil {
		return err
	}

	return "deleted"
}

func queryToDocument(id string, query url.Values) c.Document {
	document := c.NewDocument(id)

	for key, element := range query {
		document[key] = element[0]
	}

	return document
}
