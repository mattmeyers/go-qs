package qs

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/spf13/cast"
)

var (
	// ErrInvalidQS will be returned when the provided query string cannot
	// be parsed.
	ErrInvalidQS = errors.New("invalid query string")
	// ErrUnbalanced will be returned when brackets in the query string are
	// unbalanced e.g. a[[b] or a[b]].
	ErrUnbalanced = errors.New("brackets are unbalanced")
)

// QS holds both the raw query string as well as the parsed data structure.
// If a subkey is not provided, an empty string is used as the subkey. All
// operations on this struct are safe for concurrent use.
type QS struct {
	// RawQuery is the string passed to New.
	RawQuery string
	// Values holds the parsed data. The top level node is a placeholder with
	// the key "" and no values.
	Values *node
	// MaxDepth is the number of subkeys to parse before stopping (Default: 5)
	MaxDepth int
	// PathDelimiter is the string that separates keys in the path. Providing
	// a delimiter overrides the default behavior of supplying a path as
	// variadic arguments to Get, Add, Set, etc. If this option is set, then
	// any variadic methods will instead split the first parameter on this
	// value and treat the resulting slice as the path elements. (Default: "")
	PathDelimiter string

	mutex *sync.RWMutex
}

type node struct {
	Key      string
	Values   []interface{}
	Children map[string]*node
}

// Option is a functional option used to configure a new QS.
type Option func(*QS)

// MaxDepth sets the MaxDepth property of a QS struct. If MaxDepth is less
// than or equal to zero, then the query string will be parsed completely.
// Otherwise, the first MaxDepth subkeys will be parsed and the rest of the
// key will be used as the final subkey in the path e.g.
//		a[b][c][d] w/ max 2 => []string{"a", "b", "[c][d]"}
func MaxDepth(d int) Option {
	return func(qs *QS) {
		qs.MaxDepth = d
	}
}

// PathDelimiter sets the PathDelimiter property of a QS struct. Providing
// a delimiter overrides the default behavior of supplying a path as
// variadic arguments to Get, Add, Set, etc. If this option is set, then
// any variadic methods will instead split the first parameter on this
// value and treat the resulting slice as the path elements.
func PathDelimiter(d string) Option {
	return func(qs *QS) {
		qs.PathDelimiter = d
	}
}

// New produces a new QS data structure. Before parsing subkeys, the raw
// query string is processed by net/url's ParseQuery function. This function
// unescapes any URL encoding.
//
// An error will be returned if the provided query string cannot be parsed.
func New(rawQuery string, opts ...Option) (*QS, error) {
	qs := &QS{
		RawQuery: rawQuery,
		Values:   newNode(""),
		MaxDepth: 5,
		mutex:    &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(qs)
	}

	pq, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, ErrInvalidQS
	}

	if qs.PathDelimiter != "" {
		defer func(del string) { qs.PathDelimiter = del }(qs.PathDelimiter)
		qs.PathDelimiter = ""
	}

	for key, val := range pq {
		keys, err := parseKey(key, qs.MaxDepth)
		if err != nil {
			return nil, err
		}

		qs.Set(toISlice(val), keys...)
	}

	return qs, nil
}

func parseKey(key string, maxDepth int) ([]string, error) {
	inBrackets := false
	cur := make([]rune, 0)
	keys := make([]string, 0)
	depth := 0

	for i, c := range key {
		if c == '[' && !inBrackets {
			inBrackets = true
			keys = append(keys, string(cur))
			cur = cur[:0]
			depth++

			if depth == maxDepth {
				cur = []rune(key)[i:]
				break
			}
		} else if c == ']' && inBrackets {
			inBrackets = false
		} else if (c == '[' && inBrackets) || (c == ']' && !inBrackets) {
			return nil, ErrInvalidQS
		} else {
			cur = append(cur, c)
		}
	}

	keys = append(keys, string(cur))

	return keys, nil
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
	pLen := len(path)

	if path[pLen-1] == "" {
		path = path[0 : pLen-1]
		pLen--
	}

	currNode := q.Values
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
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.PathDelimiter != "" && len(path) > 0 {
		path = strings.Split(path[0], q.PathDelimiter)
	}

	n := q.navigate(path...)
	if n != nil {
		n.Values = vals
	}
}

// Add follows the provided path and appends the given value
// to the list of values at the end.
func (q *QS) Add(val interface{}, path ...string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.PathDelimiter != "" && len(path) > 0 {
		path = strings.Split(path[0], q.PathDelimiter)
	}

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
	if len(path) == 0 {
		return nil
	}

	q.mutex.RLock()
	defer q.mutex.RUnlock()

	currNode := q.Values
	var ok bool

	if q.PathDelimiter != "" {
		path = strings.Split(path[0], q.PathDelimiter)
	}
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

// GetStringSlice retrieves all values at a given path as a string slice.
func (q *QS) GetStringSlice(path ...string) []string {
	return cast.ToStringSlice(q.GetAll(path...))
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
	if len(path) == 0 {
		return make([]interface{}, 0)
	}

	q.mutex.RLock()
	defer q.mutex.RUnlock()

	currNode := q.Values
	var ok bool

	if q.PathDelimiter != "" {
		path = strings.Split(path[0], q.PathDelimiter)
	}
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

// String converts a QS data structure into its string form. This string is not
// properly encoded for use in URLs. Use EncodedString to retrieve an encoded
// version of the query string.
func (q *QS) String() string {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return print("", q.Values, false)
}

// EncodedString converts a QS data structure into its string form. All keys and
// values are encoded and ready for use in a URL.
func (q *QS) EncodedString() string {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return print("", q.Values, true)
}

func print(key string, n *node, encode bool) string {
	if key == "" {
		key = n.Key
	} else {
		key = fmt.Sprintf("%s[%s]", key, n.Key)
	}

	s := make([]string, 0)
	for _, val := range n.Values {
		if encode {
			s = append(s, fmt.Sprintf("%s=%v", url.QueryEscape(key), url.QueryEscape(fmt.Sprintf("%v", val))))
		} else {
			s = append(s, fmt.Sprintf("%s=%v", key, val))
		}
	}

	if len(n.Children) == 0 {
		return strings.Join(s, "&")
	}

	for _, v := range n.Children {
		s = append(s, print(key, v, encode))
	}

	return strings.Join(s, "&")
}
