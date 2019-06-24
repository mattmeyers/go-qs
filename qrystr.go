package qrystr

import (
	"net/url"
	"regexp"
	"strings"
)

// QS holds both the raw query string as well as the parsed data structure.
// If a subkey is not provided, an empty string is used as the subkey.
type QS struct {
	RawQuery string
	Values   map[string]map[string][]string
}

// NewQS produces a new QS data structure. Before parsing subkeys, the raw
// query string is processed by net/url's ParseQuery function. This function
// unescapes and URL encoding.
func NewQS(rawQuery string) (*QS, error) {
	qs := &QS{RawQuery: rawQuery, Values: make(map[string]map[string][]string)}
	pq, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile(`([a-zA-Z0-9]+)(?:\[([a-zA-Z0-9]+)?\])?`)
	if err != nil {
		return nil, err
	}
	for key, val := range pq {
		matches := r.FindStringSubmatch(key)
		if qs.Values[matches[1]] == nil {
			qs.Values[matches[1]] = make(map[string][]string)
		}
		qs.Values[matches[1]][matches[2]] = append(qs.Values[matches[1]][matches[2]], val...)
	}

	return qs, nil
}

// Get retrieves the value at a comma separated path in the structure of
// values. If there are multiple values at the provided path, only the first
// is returned. Use GetAll to retrieve all values. This function does not
// return an error. If a value is not found, then the empty string is returned.
func (q *QS) Get(path string) string {
	parts := getPathComponents(path)
	if val, ok := q.Values[parts[0]][parts[1]]; ok {
		return val[0]
	}
	return ""
}

// GetAll retrieves all values at a given comma separated path in the values
// structure. No error is returned from this function. If no values exists at
// the given path, then an empty string slice is returned.
func (q *QS) GetAll(path string) []string {
	parts := getPathComponents(path)
	if val, ok := q.Values[parts[0]][parts[1]]; ok && len(val) != 0 {
		return val
	}
	return []string{}
}

func getPathComponents(path string) []string {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		parts = append(parts, "")
	}
	return parts
}
