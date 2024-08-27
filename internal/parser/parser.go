package parser

import (
	"encoding/json"
	"io"
)

// Parse parses the response from Elasticsearch with profile enabled.
func Parse(r io.Reader, method, path, query string) (*Response, error) {
	response := Response{
		Method: method,
		Path:   path,
		Query:  query,
	}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
