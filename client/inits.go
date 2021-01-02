package client

import (
	"encoding/json"
	"io"
	"net/http"
)

// Unmarshaller wraps json.Unmarshal
var Unmarshaller func(data []byte, v interface{}) error

// Marshaller wraps json.Marshall
var Marshaller func(v interface{}) ([]byte, error)

// RequestCreator wraps http.NewRequest
var RequestCreator func(method, url string, body io.Reader) (*http.Request, error)

func init() {
	Marshaller = json.Marshal
	Unmarshaller = json.Unmarshal
	RequestCreator = http.NewRequest
}
