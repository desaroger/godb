package http

import (
	"encoding/json"
	"errors"
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
	http.HandleFunc("/", api.main)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (api *httpJsonApi) main(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	var response any
	if path == "_set" {
		document := queryToDocument(r.URL.Query())
		response = api.set(document)
	} else if path == "_patch" {
		document := queryToDocument(r.URL.Query())
		response = api.set(document)
	} else if strings.HasSuffix(path, "_list") {
		id := strings.TrimSuffix(path, "_list")
		id = strings.Trim(path, "/")
		response = api.list(id)
	} else if path != "" {
		response = api.get(path)
	} else {
		response = api.home()
	}

	if err, ok := response.(error); ok {
		response = c.Document{
			"error": err.Error(),
		}
	}

	json.NewEncoder(w).Encode(response)
}

func (api *httpJsonApi) home() any {
	ids, err := api.godb.List("")
	if err != nil {
		return err
	}

	return ids
}

func (api *httpJsonApi) get(id string) any {
	document, err := api.godb.Get(id)
	if err != nil {
		if errors.Is(err, c.ErrDocumentDoestNotExist) {
			var ids []string
			ids, err = api.godb.List(id)
			if err != nil {
				return err
			}
			return ids
		}
		return err
	}

	return document
}

func (api *httpJsonApi) set(document c.Document) any {
	err := api.godb.Set(document)
	if err != nil {
		return err
	}

	return true
}

func (api *httpJsonApi) patch(document c.Document) any {
	err := api.godb.Patch(document)
	if err != nil {
		return err
	}

	return true
}

func (api *httpJsonApi) list(id string) any {
	ids, err := api.godb.List(id)
	if err != nil {
		return err
	}

	return ids
}

func apiSuccess(w http.ResponseWriter, v any) {
	json.NewEncoder(w).Encode(v)
}

func apiError(w http.ResponseWriter, err error) {
	response := map[string]string{
		"error": err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}

func queryToDocument(query url.Values) c.Document {
	document := c.NewDocument("")

	for key, element := range query {
		document[key] = element[0]
	}

	return document
}
