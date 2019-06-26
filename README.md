# qrystr

This is a query string parsing library. It uses `net/url`'s `ParseQuery` function as a starting point, then parses additional subkeys from the returned keys. In order to be properly parsed, a query string should take the form `key1=value1&key2[]=value2&key3[subkey3]=value3...`. This query string will be parsed into the structure
```json
{
  "key1": {
    "": [
      "value1"
    ],
  },
  "key2": {
    "": [
      "value2"
    ],
  },
  "key3": {
    "subkey3": [
      "value3"
    ],
  },
}
```

If mutiple values are provided for the same key or subkey, the values are placed in a slice. Note that missing subkeys will use the empty string as a key. Additionally, a single key can have multiple subkeys. For example, the string 
```
k1[sk1]=v1&k1[sk1]=v2&k1[sk2]=v3&k2[]=v4&k2[]=v5&k3=v6&k3=v7
```
will parse to 
```json
{
  "k1": {
    "sk1": [
      "v1", "v2"
    ],
    "sk2": [
      "v3"
    ]
  },
  "k2": {
    "": [
      "v4", "v5"
    ]
  },
  "k3": {
    "": [
      "v6", "v7"
    ]
  },
}
```

## Initializing

To parse a query string, use the `qrystr.NewQS(rawQuery string)` function. This will parse the query string and return the `QS` struct containing the values data structure. If any errors occur while parsing, the error will be returned.

## Retrieving Values

The values data structure can always be directly accessed in the `QS.Values` property. Alternatively, the keys can be provided to the variadic `Get` or `GetAll` methods. If multiple values exist at a given path, `Get` will return the first in the slice while `GetAll` will return the entire slice. Neither of these methods return an error. If no values exist at the given path, then either `""` or `[]string{}` is returned.

```Go
package main

import (
	"fmt"

	"bitbucket.org/mcgstrategic/qrystr"
)

func main() {
	query := "k1[sk1]=v1&k1[sk1]=v2&k1[sk2]=v3&k2[]=v4&k2[]=v5&k3=v6&k3=v7"
	qs, _ := qrystr.NewQS(query)

	fmt.Printf("%v\n", qs.Values)
	// map[k1:map[sk1:[v1 v2] sk2:[v3]] k2:map[:[v4 v5]] k3:map[:[v6 v7]]]
	fmt.Printf("%v\n", qs.Get("k1", "sk1"))
	// v1
	fmt.Printf("%v\n", qs.GetAll("k1", "sk1"))
	// [v1 v2]
	fmt.Printf("%v\n", qs.Get("k2"))
	// v4
	fmt.Printf("%v\n", qs.GetAll("k3"))
	// [v6 v7]
	fmt.Printf("%v\n", qs.Get("k1", "sk1"))
	// v1
	fmt.Printf("%v\n", qs.Get("k1", "sk3"))
	// ""
	fmt.Printf("%v\n", qs.GetAll("k1", "sk3"))
	// []
}
```