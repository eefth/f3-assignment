package client

import "encoding/json"

// Unmarshaller overrides the go build-in json.Unmarshal
var Unmarshaller func(data []byte, v interface{}) error

// Marshaller overrides build in json.Marshall
var Marshaller func(v interface{}) ([]byte, error)

func init() {
	Marshaller = json.Marshal
	Unmarshaller = json.Unmarshal
}
