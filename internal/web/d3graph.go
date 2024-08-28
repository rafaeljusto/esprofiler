package web

import (
	"time"

	"github.com/rafaeljusto/esprofiler/internal/parser"
)

type d3Data struct {
	Name     string   `json:"name"`
	Value    int64    `json:"value"`
	Children []d3Data `json:"children"`
}

func buildSearchD3Data(search []parser.SearchQuery) []d3Data {
	var result []d3Data
	for _, s := range search {
		result = append(result, d3Data{
			Name:     s.Type + " - " + s.Description,
			Value:    time.Duration(s.Took).Milliseconds(),
			Children: buildSearchD3Data(s.Children),
		})
	}
	return result
}

func buildCollectorD3Data(collector []parser.SearchCollector) []d3Data {
	var result []d3Data
	for _, c := range collector {
		result = append(result, d3Data{
			Name:     c.Name + " - " + c.Reason,
			Value:    time.Duration(c.Took).Milliseconds(),
			Children: buildCollectorD3Data(c.Children),
		})
	}
	return result
}
