package riot

import "strings"

// Region denotes various regions of League of Legends
type Region string

const (
	// NA1 is North America 1
	NA1 Region = "NA1"
	// TODO
)

// ParseRegion parses string into a region
func ParseRegion(s string) Region {
	return Region(strings.ToUpper(s))
}

// ToLower converts region to bigquery-compatible lowercase
func (r Region) ToLower() string {
	return strings.ToLower(string(r))
}
