package qs

import (
	"net/url"
	"strings"

	"github.com/spf13/cast"
)

// QS holds both the raw query string as well as the parsed data structure.
// If a subkey is not provided, an empty string is used as the subkey.
type QS struct {
	RawQuery string
	Values   *node
}

type node struct {
	Key      string
	Values   []interface{}
	Children map[string]*node
}

// New produces a new QS data structure. Before parsing subkeys, the raw
// query string is processed by net/url's ParseQuery function. This function
// unescapes and URL encoding.
func New(rawQuery string) (*QS, error) {
	qs := &QS{RawQuery: rawQuery, Values: newNode("")}
	pq, err := url.ParseQuery(rawQuery)

	if err != nil {
		return nil, err
	}

	for key, val := range pq {
		keys := strings.Split(key, "[")
		for i, k := range keys {
			keys[i] = strings.Trim(k, "]")
		}

		qs.Set(toISlice(val), keys...)
	}

	return qs, nil
}

func toISlice(s []string) []interface{} {
	iSlice := make([]interface{}, len(s))
	for i, e := range s {
		iSlice[i] = interface{}(e)
	}
	return iSlice
}

func newNode(key string) *node {
	return &node{
		Key:      key,
		Values:   make([]interface{}, 0),
		Children: make(map[string]*node),
	}
}

func (q *QS) navigate(path ...string) *node {
	currNode := q.Values

	pLen := len(path)

	if path[pLen-1] == "" {
		path = path[0 : pLen-1]
		pLen--
	}

	for i, p := range path {

		childNode, ok := currNode.Children[p]
		if !ok {
			currNode.Children[p] = newNode(p)
			childNode = currNode.Children[p]
		}
		currNode = childNode

		if i+1 == pLen {
			return currNode
		}
	}

	return nil
}

// Set follows the provided path and overwrites the values
// at the end with the provided values.
func (q *QS) Set(vals []interface{}, path ...string) {
	n := q.navigate(path...)
	if n != nil {
		n.Values = vals
	}
}

// Add follows the provided path and appends the given value
// to the list of values at the end.
func (q *QS) Add(val interface{}, path ...string) {
	n := q.navigate(path...)
	if n != nil {
		n.Values = append(n.Values, val)
	}
}

// Get follows the provided keys and returns the value at the end.
// If there are multiple values at the provided path, only the first
// is returned. Use GetAll to retrieve all values. This function does not
// return an error. If a value is not found, then nil is returned.
func (q *QS) Get(path ...string) interface{} {
	currNode := q.Values
	var ok bool

	pLen := len(path)

	for i, p := range path {
		currNode, ok = currNode.Children[p]
		if !ok {
			return nil
		}

		if i+1 == pLen {
			if len(currNode.Values) > 0 {
				return currNode.Values[0]
			}
			return nil
		}
	}

	return nil
}

// GetString retrieves the value at the given path as a string.
func (q *QS) GetString(path ...string) string {
	return cast.ToString(q.Get(path...))
}

// GetInt retrieves the value at the given path as an int. If
// the value cannot be converted to an int, 0 is returned.
func (q *QS) GetInt(path ...string) int {
	return cast.ToInt(q.Get(path...))
}

// GetInt32 retrieves the value at the given path as an int32. If
// the value cannot be converted to an int, 0 is returned.
func (q *QS) GetInt32(path ...string) int32 {
	return cast.ToInt32(q.Get(path...))
}

// GetInt64 retrieves the value at the given path as an int64. If
// the value cannot be converted to an int, 0 is returned.
func (q *QS) GetInt64(path ...string) int64 {
	return cast.ToInt64(q.Get(path...))
}

// GetFloat32 retrieves the value at the given path as a float32. If
// the value cannot be converted to a float, 0 is returned.
func (q *QS) GetFloat32(path ...string) float32 {
	return cast.ToFloat32(q.Get(path...))
}

// GetFloat64 retrieves the value at the given path as a float64. If
// the value cannot be converted to a float, 0 is returned.
func (q *QS) GetFloat64(path ...string) float64 {
	return cast.ToFloat64(q.Get(path...))
}

// GetBool retrieves the value at the given path as a bool. If
// the value cannot be converted to a float, false is returned.
func (q *QS) GetBool(path ...string) bool {
	return cast.ToBool(q.Get(path...))
}

// GetWithDefault follows the provided keys and returns the value at the end.
// If there are multiple values at the provided path, only the first
// is returned. Use GetAllWithDefault to retrieve all values. This function does not
// return an error. If a value is not found, then the provided default is returned.
func (q *QS) GetWithDefault(def interface{}, path ...string) interface{} {
	val := q.Get(path...)
	if val == nil {
		return def
	}
	return val
}

// GetAll follows the provided keys and returns all values at the end.
// No error is returned from this function. If no values exists at
// the given path, then a slice of interfaces is returned.
func (q *QS) GetAll(path ...string) []interface{} {
	currNode := q.Values
	var ok bool

	pLen := len(path)

	for i, p := range path {
		currNode, ok = currNode.Children[p]
		if !ok {
			return nil
		}

		if i+1 == pLen {
			return currNode.Values
		}
	}

	return nil
}

// GetAllWithDefault follows the provided keys and returns all values at the end.
// No error is returned from this function. If no values exists at
// the given path, then the provided default value is returned.
func (q *QS) GetAllWithDefault(def []interface{}, path ...string) []interface{} {
	vals := q.GetAll(path...)
	if len(vals) == 0 {
		return def
	}
	return vals
}
