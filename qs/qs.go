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

func (q *QS) Set(vals []interface{}, path ...string) {
	n := q.navigate(path...)
	if n != nil {
		n.Values = vals
	}
}

func (q *QS) Add(val interface{}, path ...string) {
	n := q.navigate(path...)
	if n != nil {
		n.Values = append(n.Values, val)
	}
}

// Get follows the provided keys and returns the value at the end.
// If there are multiple values at the provided path, only the first
// is returned. Use GetAll to retrieve all values. This function does not
// return an error. If a value is not found, then the empty string is returned.
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
			return currNode.Values[0]
		}
	}

	return nil
}

func (q *QS) GetString(path ...string) string {
	return cast.ToString(q.Get(path...))
}

func (q *QS) GetInt(path ...string) int {
	return cast.ToInt(q.Get(path...))
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
// the given path, then an empty string slice is returned.
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
