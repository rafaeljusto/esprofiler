package parser

import (
	"encoding/json"
	"io"
)

// Parse parses the response from Elasticsearch with profile enabled.
func Parse(r io.Reader) (*Response, error) {
	var response Response
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
