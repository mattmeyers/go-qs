# qs

This is a query string parsing library inspired by [qs](https://www.npmjs.com/package/qs). It uses `net/url`'s `ParseQuery` function as a starting point, then parses additional subkeys from the returned keys. 

# Installation

This library supports modules, or you can install with `go get`.
```
go get -u github.com/mattmeyers/go-qs/qs
```

# Usage

## Parsing

Given a query string such as `a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true`, first initialize a QS struct with 

```go
qs.New("a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true")
```

This function will return an error if there was a problem parsing the query string. 

## Getting Values

After parsing, do not try to navigate the tree structure manually, but rather use one of the many provided getter methods. To provide flexibility, values are stored in a slice of interfaces. There are four generic getters to get these values without and type conversions.

- `Get(path ...string) interface{}`
- `GetWithDefault(def interface{}, path ...string) interface{}`
- `Get(path ...string) []interface{}`
- `GetWithDefault(def []interface{}, path ...string) []interface{}`

This library also provides getters for specific data types using the [cast](https://github.com/spf13/cast) library. If any type conversions fail, the types zero value is returned.

- `GetString(path ...string) string`
- `GetInt(path ...string) int`
- `GetInt32(path ...string) int32`
- `GetInt64(path ...string) int64`
- `GetFloat32(path ...string) float32`
- `GetFloat64(path ...string) float64`
- `GetBool(path ...string) bool`

For example, suppose we have the query string `a[b]=3&c[d][e]=true`. We can retrieve both values using

```go
q, _ := qs.New("a[b]=3&c[d][e]=true")

firstVal := q.GetInt("a", "b")
// firstVal == 3

secondVal := q.getBool("c", "d", "e")
// secondVal == true
```

## Setting Values

Provided a parsed query string, values can be set or added. 

- `Set(values []interface{}, path ...string)`
- `Add(value interface{}, path ...string)`

Note that setting completely overwrites the previous values while adding simply appends a new value to the list.  For example

```go
q, _ := qs.New("a[b]=3&c[d][e]=true")

q.Set([]interface{1, "a", true}, "a", "b")
firstVals := q.GetAll("a", "b")
// firstVals == []interface{1, "a", true}

q.Add("def", "c", "d", "e")
secondVals := q.GetAll("c", "d", "e")
// firstVals == []interface{true, "def"}
```

# TODO

- Add `ToMap() map[string]interface{}` method
- Add `ToString() string` method
- Add parsing options
