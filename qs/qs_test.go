package qs

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewQS(t *testing.T) {
	query := "a[b]=c" //a[b][c][d]=c&a[g]=h&a[g]=i&d[]=f&j=k"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	// q.Set([]interface{}{"d"}, "a", "b")
	// out, _ := json.Marshal(q.Values.Children)
	fmt.Printf("%T - %v\n", q.Get("a", "b"), q.Get("a", "b"))
	fmt.Printf("%T - %v\n", q.GetString("a", "b"), q.GetString("a", "b"))
	fmt.Printf("%T - %v\n", q.GetInt("a", "b"), q.GetInt("a", "b"))

	exp := &node{
		Key:    "",
		Values: make([]interface{}, 0),
		Children: map[string]*node{
			"a": &node{
				Key:    "a",
				Values: make([]interface{}, 0),
				Children: map[string]*node{
					"b": {
						Key:      "b",
						Values:   []interface{}{"c"},
						Children: make(map[string]*node),
					},
				},
			},
		}}

	if !(cmp.Equal(q.Values, exp)) {
		t.Fatalf("Generated incorrect tree: %s", cmp.Diff(q.Values, exp))
	}
}

// func TestGet(t *testing.T) {
// 	query := "a[b]=c&a[g]=h&a[g]=i&d[]=f&j=k&w->x[y]=z&q.u=v&m[n]=o:p"
// 	q, _ := qs.New(query)

// 	table := []struct {
// 		Desc   string
// 		Key    string
// 		Subkey string
// 		Exp    string
// 	}{
// 		{"a[b]", "a", "b", "c"},
// 		{"a[g]", "a", "g", "h"},
// 		{"d[]", "d", "", "f"},
// 		{"j", "j", "", "k"},
// 		{"FAKE KEY", "foo", "bar", ""},
// 		{"No key", "", "", ""},
// 		{"w->x", "w->x", "y", "z"},
// 		{"q.u", "q.u", "", "v"},
// 		{"m[n]", "m", "n", "o:p"},
// 	}

// 	for _, row := range table {
// 		t.Run(row.Desc, func(t *testing.T) {
// 			path := []string{row.Key}
// 			if row.Subkey != "" {
// 				path = append(path, row.Subkey)
// 			}
// 			v := q.Get(path...)
// 			if v != row.Exp {
// 				t.Fatalf("Get() failed: expected %s, got %s", row.Exp, v)
// 			}
// 		})
// 	}
// }

// func TestGetAll(t *testing.T) {
// 	query := "a[b]=c&a[g]=h&a[g]=i&d[]=f&j=k&w->x[y]=z&q.u=v&m[n]=o:p"
// 	q, _ := qs.NewQS(query)
// 	table := []struct {
// 		Desc   string
// 		Key    string
// 		Subkey string
// 		Exp    []string
// 	}{
// 		{"a[b]", "a", "b", []string{"c"}},
// 		{"a[g]", "a", "g", []string{"h", "i"}},
// 		{"d[]", "d", "", []string{"f"}},
// 		{"j", "j", "", []string{"k"}},
// 		{"FAKE KEY", "foo", "bar", []string{}},
// 		{"No key", "", "", []string{}},
// 		{"w->x", "w->x", "y", []string{"z"}},
// 		{"q.u", "q.u", "", []string{"v"}},
// 		{"m[n]", "m", "n", []string{"o:p"}},
// 	}

// 	for _, row := range table {
// 		t.Run(row.Desc, func(t *testing.T) {
// 			path := []string{row.Key}
// 			if row.Subkey != "" {
// 				path = append(path, row.Subkey)
// 			}
// 			v := q.GetAll(path...)
// 			if !cmp.Equal(v, row.Exp) {
// 				t.Fatalf("Get() failed: expected %v, got %v", row.Exp, v)
// 			}
// 		})
// 	}
// }
