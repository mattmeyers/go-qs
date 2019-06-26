package qrystr_test

import (
	"testing"

	"bitbucket.org/mcgstrategic/qrystr"
	"github.com/google/go-cmp/cmp"
)

func TestNewQS(t *testing.T) {
	query := "a[b]=c&a[g]=h&a[g]=i&d[]=f&j=k"
	q, err := qrystr.NewQS(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}
	exp := map[string]map[string][]string{
		"a": map[string][]string{"b": []string{"c"}, "g": []string{"h", "i"}},
		"d": map[string][]string{"": []string{"f"}},
		"j": map[string][]string{"": []string{"k"}},
	}
	if !(cmp.Equal(q.Values, exp)) {
		t.Fatalf("Generated value incorrect: expected %v, got %v", exp, q.Values)
	}
}

func TestGet(t *testing.T) {
	query := "a[b]=c&a[g]=h&a[g]=i&d[]=f&j=k&w->x[y]=z&q.u=v"
	q, _ := qrystr.NewQS(query)

	table := []struct {
		Desc   string
		Key    string
		Subkey string
		Exp    string
	}{
		{"a[b]", "a", "b", "c"},
		{"a[g]", "a", "g", "h"},
		{"d[]", "d", "", "f"},
		{"j", "j", "", "k"},
		{"FAKE KEY", "m", "n", ""},
		{"No key", "", "", ""},
		{"w->x", "w->x", "y", "z"},
		{"q.u", "q.u", "", "v"},
	}

	for _, row := range table {
		t.Run(row.Desc, func(t *testing.T) {
			path := []string{row.Key}
			if row.Subkey != "" {
				path = append(path, row.Subkey)
			}
			v := q.Get(path...)
			if v != row.Exp {
				t.Fatalf("Get() failed: expected %s, got %s", row.Exp, v)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	query := "a[b]=c&a[g]=h&a[g]=i&d[]=f&j=k&w->x[y]=z&q.u=v"
	q, _ := qrystr.NewQS(query)
	table := []struct {
		Desc   string
		Key    string
		Subkey string
		Exp    []string
	}{
		{"a[b]", "a", "b", []string{"c"}},
		{"a[g]", "a", "g", []string{"h", "i"}},
		{"d[]", "d", "", []string{"f"}},
		{"j", "j", "", []string{"k"}},
		{"FAKE KEY", "m", "n", []string{}},
		{"No key", "", "", []string{}},
		{"w->x", "w->x", "y", []string{"z"}},
		{"q.u", "q.u", "", []string{"v"}},
	}

	for _, row := range table {
		t.Run(row.Desc, func(t *testing.T) {
			path := []string{row.Key}
			if row.Subkey != "" {
				path = append(path, row.Subkey)
			}
			v := q.GetAll(path...)
			if !cmp.Equal(v, row.Exp) {
				t.Fatalf("Get() failed: expected %v, got %v", row.Exp, v)
			}
		})
	}
}
