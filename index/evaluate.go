package index

import (
	"encoding/json"

	v8 "rogchap.com/v8go"

	c "godb/common"
)

func evaluate(document c.Document, index_func string) (string, c.Document, error) {
	document_bytes, err := json.Marshal(document)
	if err != nil {
		return "", nil, err
	}
	document_json := string(document_bytes)

	iso := v8.NewIsolate()
	ctx := v8.NewContext(iso)
	_, err = ctx.RunScript(c.S("const indexFunc = %s", index_func), "main.js")
	if err != nil {
		return "", nil, err
	}
	_, err = ctx.RunScript(c.S("const indexResult = indexFunc(%s)", document_json), "main.js")
	if err != nil {
		return "", nil, err
	}
	indexId, err := ctx.RunScript("indexResult && indexResult[0] || null", "main.js")
	if err != nil {
		return "", nil, err
	}
	if indexId.IsNull() {
		return "", nil, nil
	}
	indexContent, err := ctx.RunScript("JSON.stringify(indexResult && indexResult[1] || null)", "main.js")
	if err != nil {
		return "", nil, err
	}

	var result_document c.Document
	err = json.Unmarshal([]byte(indexContent.String()), &result_document)
	if err != nil {
		return "", nil, err
	}

	return indexId.String(), result_document, nil
}
