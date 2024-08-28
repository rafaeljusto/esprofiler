package parser

import (
	"regexp"
	"strconv"
	"time"
)

var reShardID = regexp.MustCompile(`^\[.+\]\[(?P<index>.+)\]\[(?P<shard>\d+)\]$`)

// Response represents the response from Elasticsearch with profile enabled.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-profile.html
type Response struct {
	Took    DurationInMilliseconds `json:"took"`
	Profile struct {
		Shards Shards `json:"shards"`
	} `json:"profile"`
}

// Shard represents a single shard in the profile.
type Shard struct {
	ID       ShardID  `json:"id"`
	Searches []Search `json:"searches"`
}

// Shards is a collection of shards.
type Shards []Shard

// GroupByIndex groups the shards by index.
func (s Shards) GroupByIndex() map[string][]Shard {
	grouped := make(map[string][]Shard)
	for _, shard := range s {
		grouped[shard.ID.Index] = append(grouped[shard.ID.Index], shard)
	}
	return grouped
}

// ShardID contains all information inside the shard ID.
type ShardID struct {
	Index string
	Shard int64
}

// UnmarshalJSON unmarshals the shard ID.
func (s *ShardID) UnmarshalJSON(data []byte) error {
	raw, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	match := reShardID.FindStringSubmatch(raw)
	result := make(map[string]string)
	for i, name := range reShardID.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	s.Index = result["index"]
	s.Shard, err = strconv.ParseInt(result["shard"], 10, 64)
	if err != nil {
		return err
	}
	return nil
}

// Search represents a single search in the profile.
type Search struct {
	Query       []SearchQuery         `json:"query"`
	RewriteTime DurationInNanoseconds `json:"rewrite_time"`
	Collector   []SearchCollector     `json:"collector"`
}

// Took returns the total time for the search.
func (p Search) Took() time.Duration {
	var total time.Duration
	for _, query := range p.Query {
		total += time.Duration(query.Took)
	}
	for _, collector := range p.Collector {
		total += time.Duration(collector.Took)
	}
	return total + time.Duration(p.RewriteTime)
}

// SearchQuery represents a single search query in the profile.
type SearchQuery struct {
	Type        string                `json:"type"`
	Description string                `json:"description"`
	Took        DurationInNanoseconds `json:"time_in_nanos"`
	Breakdown   map[string]int64      `json:"breakdown"`
	Children    []SearchQuery         `json:"children"`
}

// SearchCollector represents a single search collector in the profile.
type SearchCollector struct {
	Name     string                `json:"name"`
	Reason   string                `json:"reason"`
	Took     DurationInNanoseconds `json:"time_in_nanos"`
	Children []SearchCollector     `json:"children"`
}

// DurationInMilliseconds represents a time.Duration in milliseconds.
type DurationInMilliseconds time.Duration

// UnmarshalJSON unmarshals the time.Duration in milliseconds.
func (d *DurationInMilliseconds) UnmarshalJSON(data []byte) error {
	number, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*d = DurationInMilliseconds(time.Duration(number) * time.Millisecond)
	return nil
}

// String returns the human readable representation.
func (d DurationInMilliseconds) String() string {
	return time.Duration(d).String()
}

// DurationInNanoseconds represents a time.Duration in nanoseconds.
type DurationInNanoseconds time.Duration

// UnmarshalJSON unmarshals the time.Duration in nanoseconds.
func (d *DurationInNanoseconds) UnmarshalJSON(data []byte) error {
	number, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*d = DurationInNanoseconds(time.Duration(number))
	return nil
}

// String returns the human readable representation.
func (d DurationInNanoseconds) String() string {
	return time.Duration(d).String()
}
