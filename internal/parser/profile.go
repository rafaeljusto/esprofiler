package parser

import (
	"crypto/md5"
	"fmt"
	"time"
)

// Response represents the response from Elasticsearch with profile enabled.
type Response struct {
	Took    int64 `json:"took"`
	Profile struct {
		Shards []struct {
			ID       string `json:"id"`
			Searches []struct {
				Query       []ProfileSearchQuery     `json:"query"`
				RewriteTime int64                    `json:"rewrite_time"`
				Collector   []ProfileSearchCollector `json:"collector"`
			} `json:"searches"`
		} `json:"shards"`
	} `json:"profile"`
}

// ProfileSearchQuery represents a single search query in the profile.
type ProfileSearchQuery struct {
	Type        string               `json:"type"`
	Description string               `json:"description"`
	TimeInNanos int64                `json:"time_in_nanos"`
	Breakdown   map[string]int64     `json:"breakdown"`
	Children    []ProfileSearchQuery `json:"children"`
}

// HexColor returns the hex color for the specific search query.
func (p ProfileSearchQuery) HexColor() string {
	encoder := md5.New()
	encoder.Write([]byte(p.Type))
	encoder.Write([]byte(p.Description))
	encoder.Write([]byte(fmt.Sprintf("%d", p.TimeInNanos)))
	for _, child := range p.Children {
		encoder.Write([]byte(child.HexColor()))
	}
	return fmt.Sprintf("#%x", encoder.Sum(nil)[:3])
}

// Took returns the total time in an human readable format.
func (p ProfileSearchQuery) Took() string {
	return time.Duration(p.TimeInNanos).String()
}

// ProfileSearchCollector represents a single search collector in the profile.
type ProfileSearchCollector struct {
	Name        string                   `json:"name"`
	Reason      string                   `json:"reason"`
	TimeInNanos int64                    `json:"time_in_nanos"`
	Children    []ProfileSearchCollector `json:"children"`
}

// HexColor returns the hex color for the specific search collector.
func (p ProfileSearchCollector) HexColor() string {
	encoder := md5.New()
	encoder.Write([]byte(p.Name))
	encoder.Write([]byte(p.Reason))
	encoder.Write([]byte(fmt.Sprintf("%d", p.TimeInNanos)))
	for _, child := range p.Children {
		encoder.Write([]byte(child.HexColor()))
	}
	return fmt.Sprintf("#%x", encoder.Sum(nil)[:3])
}

// Took returns the total time in an human readable format.
func (p ProfileSearchCollector) Took() string {
	return time.Duration(p.TimeInNanos).String()
}
