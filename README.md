# qs

This is a query string parsing library inspired by nodeJS's [qs](https://www.npmjs.com/package/qs) library. It uses `net/url`'s `ParseQuery` function as a starting point, then parses additional subkeys to create a tree of keys and values. Conversely, a query string can be programmatically created with this library. All functionality provided by this library is safe for concurrent use.

# Installation

```
go get -u github.com/mattmeyers/go-qs/qs
```

# Usage

## Parsing

A QS struct can be initialized using the
```
New(rawQuery string, opts ...Option) (*QS, error)
```
function. An empty string can be passed if the user wants to build a new query string. Functional options can be passed through this function to configure the new QS. At this time, the available options are:

* `MaxDepth(d int)` - Sets the max number of subkeys that will be parsed before stopping. Pass a non positive integer to parse all subkeys regardless of depth. Defaults to 5.

This function will return one of two errors if parsing fails

* `qs.ErrInvalidQS` - Returned when `net/url.ParseQuery` fails to parse the provided query string.
* `qs.ErrUnbalanced` - Returned when the provided query string has a key with unbalanced brackets e.g. `a[[b]=2`.

For example, given a query string such as
```
a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true
```
a QS struct with a max depth of 3 can be initialized with

```go
qs.New(
  "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true",
  qs.MaxDepth(3),
)
```

## Getting Values

After parsing, do not try to navigate the tree structure manually, but rather use one of the many provided getter methods. To provide flexibility, values are stored in a slice of interfaces. There are four generic getters to get these values without any type conversions.

- `Get(path ...string) interface{}`
- `GetWithDefault(def interface{}, path ...string) interface{}`
- `GetAll(path ...string) []interface{}`
- `GetAllWithDefault(def []interface{}, path ...string) []interface{}`

This library also provides getters for specific data types using the [cast](https://github.com/spf13/cast) library. If any type conversions fail, the type's zero value is returned.

- `GetString(path ...string) string`
- `GetStringSlice(path ...string) []string`
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

Note that setting completely overwrites the previous values while adding simply appends a new value to the list. If no node exists at the provided path, then a new node is created. For example

```go
q, _ := qs.New("a[b]=3&c[d][e]=true")

q.Set([]interface{1, "a", true}, "a", "b")
firstVals := q.GetAll("a", "b")
// firstVals == []interface{1, "a", true}

q.Add("def", "c", "d", "e")
secondVals := q.GetAll("c", "d", "e")
// secondVals == []interface{true, "def"}

q.Set(1, "x", "y")
thirdVals := q.GetInt("x", "y")
// thirdVals == 1
```

## Stringifying

`qs` provides two methods for converting a QS struct back into a string:

* `String() string` - Returns the string form of the QS data structure.
* `EncodedString() string` - Returns the string form of the QS data structure with all keys and values encoded for use in a URL.

# License

MIT License

Copyright (c) 2019 Matthew Meyers

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
