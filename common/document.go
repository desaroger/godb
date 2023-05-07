package common

type Document map[string]interface{}

func NewDocument(id string, args ...any) Document {
	document := make(Document)
	document["id"] = id
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		value := args[i+1]
		document[key.(string)] = value
	}

	return document
}

func (document Document) Get(key string, defaultValue interface{}) interface{} {
	value := document[key]
	if value == nil {
		return defaultValue
	}

	return value
}

func (document Document) GetId() (string, error) {
	if document == nil {
		return "", ErrEmptyDocument
	}
	if document["id"] == nil {
		return "", ErrInvalidId
	}

	return document["id"].(string), nil
}

func (document Document) GetIdOrNil() string {
	id, err := document.GetId()
	if err != nil {
		return "<?>"
	}

	return id
}

func (document Document) Patch(other Document) {
	for key, value := range other {
		document[key] = value
	}
}
