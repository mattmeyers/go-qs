package qs

import (
	"net/url"
	"regexp"
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

	r, err := regexp.Compile(`([\w-<>\.]+)(?:\[([\w-<>\.]+)?\])?`)
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

// Get follows the provided keys and returns the value at the end.
// If there are multiple values at the provided path, only the first
// is returned. Use GetAll to retrieve all values. This function does not
// return an error. If a value is not found, then the empty string is returned.
func (q *QS) Get(path ...string) string {
	path = getPathComponents(path)
	if val, ok := q.Values[path[0]][path[1]]; ok {
		return val[0]
	}
	return ""
}

// GetWithDefault follows the provided keys and returns the value at the end.
// If there are multiple values at the provided path, only the first
// is returned. Use GetAllWithDefault to retrieve all values. This function does not
// return an error. If a value is not found, then the provided default is returned.
func (q *QS) GetWithDefault(def string, path ...string) string {
	val := q.Get(path...)
	if val == "" {
		return def
	}
	return val
}

// GetAll follows the provided keys and returns all values at the end.
// No error is returned from this function. If no values exists at
// the given path, then an empty string slice is returned.
func (q *QS) GetAll(path ...string) []string {
	path = getPathComponents(path)
	if val, ok := q.Values[path[0]][path[1]]; ok && len(val) != 0 {
		return val
	}
	return []string{}
}

// GetAllWithDefault follows the provided keys and returns all values at the end.
// No error is returned from this function. If no values exists at
// the given path, then the provided default value is returned.
func (q *QS) GetAllWithDefault(def []string, path ...string) []string {
	vals := q.GetAll(path...)
	if len(vals) == 0 {
		return def
	}
	return vals
}

func getPathComponents(path []string) []string {
	if len(path) == 1 {
		return append(path, "")
	}
	return path
}
