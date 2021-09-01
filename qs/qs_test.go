package qs

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewQS(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&d[]=2.5&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	exp := &node{
		Key:    "",
		Values: make([]interface{}, 0),
		Children: map[string]*node{
			"a": {
				Key:    "a",
				Values: make([]interface{}, 0),
				Children: map[string]*node{
					"b": {
						Key:    "b",
						Values: []interface{}{"123"},
						Children: map[string]*node{
							"c": {
								Key:    "c",
								Values: make([]interface{}, 0),
								Children: map[string]*node{
									"d": {
										Key:      "d",
										Values:   []interface{}{"c"},
										Children: make(map[string]*node),
									},
								},
							},
						},
					},
					"g": {
						Key:      "g",
						Values:   []interface{}{"h", "i"},
						Children: make(map[string]*node),
					},
				},
			},
			"d": {
				Key:      "d",
				Values:   []interface{}{"1.05", "2.5"},
				Children: make(map[string]*node),
			},
			"j": {
				Key:      "j",
				Values:   []interface{}{"true"},
				Children: make(map[string]*node),
			},
		},
	}

	if !(cmp.Equal(q.Values, exp)) {
		t.Fatalf("Generated incorrect tree: %s", cmp.Diff(q.Values, exp))
	}
}

func TestQS_Get(t *testing.T) {
	type args struct {
		path []string
	}
	tests := []struct {
		name      string
		delimiter string
		args      args
		want      interface{}
	}{
		{
			name:      "Get nothing",
			delimiter: "",
			args:      args{[]string{}},
			want:      nil,
		},
		{
			name:      "Get a->b",
			delimiter: "",
			args:      args{[]string{"a", "b"}},
			want:      "123",
		},
		{
			name:      "Get a->b->c->d",
			delimiter: "",
			args:      args{[]string{"a", "b", "c", "d"}},
			want:      "c",
		},
		{
			name:      "Get a->g",
			delimiter: "",
			args:      args{[]string{"a", "g"}},
			want:      "h",
		},
		{
			name:      "Get d",
			delimiter: "",
			args:      args{[]string{"d"}},
			want:      "1.05",
		},
		{
			name:      "Get j",
			delimiter: "",
			args:      args{[]string{"j"}},
			want:      "true",
		},
		{
			name:      "Get z",
			delimiter: "",
			args:      args{[]string{"z"}},
			want:      nil,
		},
		{
			name:      "Get a",
			delimiter: "",
			args:      args{[]string{"a"}},
			want:      nil,
		},
		{
			name:      "Single byte delimiter",
			delimiter: ".",
			args:      args{[]string{"a.b.c.d"}},
			want:      "c",
		},
		{
			name:      "Multibyte delimiter",
			delimiter: "::",
			args:      args{[]string{"a::b"}},
			want:      "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(
				"a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true",
				PathDelimiter(tt.delimiter),
			)
			if err != nil {
				t.Fatalf("NewQS failed with err, %s", err)
			}

			if got := q.Get(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_Get_PathDelimiter(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g:h]=h&a[g]=i&d[]=1.05&j=true"
	tests := []struct {
		name      string
		delimiter string
		path      []string
		want      string
	}{
		{
			name:      "No delimiter",
			delimiter: "",
			path:      []string{"a", "b"},
			want:      "123",
		},
		{
			name:      "Single byte delimiter",
			delimiter: ".",
			path:      []string{"a.b.c.d"},
			want:      "c",
		},
		{
			name:      "Multibyte delimiter",
			delimiter: "::",
			path:      []string{"a::b"},
			want:      "123",
		},
		{
			name:      "Multibyte delimiter, part of delimiter in path",
			delimiter: "::",
			path:      []string{"a::g:h"},
			want:      "h",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs, err := New(query, PathDelimiter(tt.delimiter))
			if err != nil {
				t.Fatalf("unable to create new QS: %v", err)
			}

			if got := qs.Get(tt.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_GetString(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want string
	}{
		{
			name: "Get a->b",
			qs:   q,
			args: args{[]string{"a", "b"}},
			want: "123",
		},
		{
			name: "Get a->b->c->d",
			qs:   q,
			args: args{[]string{"a", "b", "c", "d"}},
			want: "c",
		},
		{
			name: "Get a->g",
			qs:   q,
			args: args{[]string{"a", "g"}},
			want: "h",
		},
		{
			name: "Get d",
			qs:   q,
			args: args{[]string{"d"}},
			want: "1.05",
		},
		{
			name: "Get j",
			qs:   q,
			args: args{[]string{"j"}},
			want: "true",
		},
		{
			name: "Get z",
			qs:   q,
			args: args{[]string{"z"}},
			want: "",
		},
		{
			name: "Get a",
			qs:   q,
			args: args{[]string{"a"}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qs.GetString(tt.args.path...); got != tt.want {
				t.Errorf("QS.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_GetInt(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want int
	}{
		{
			name: "Get a->b",
			qs:   q,
			args: args{[]string{"a", "b"}},
			want: 123,
		},
		{
			name: "Get a->b->c->d",
			qs:   q,
			args: args{[]string{"a", "b", "c", "d"}},
			want: 0,
		},
		{
			name: "Get a->g",
			qs:   q,
			args: args{[]string{"a", "g"}},
			want: 0,
		},
		{
			name: "Get d",
			qs:   q,
			args: args{[]string{"d"}},
			want: 0,
		},
		{
			name: "Get j",
			qs:   q,
			args: args{[]string{"j"}},
			want: 0,
		},
		{
			name: "Get z",
			qs:   q,
			args: args{[]string{"z"}},
			want: 0,
		},
		{
			name: "Get a",
			qs:   q,
			args: args{[]string{"a"}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qs.GetInt(tt.args.path...); got != tt.want {
				t.Errorf("QS.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_GetBool(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want bool
	}{
		{
			name: "Get a->b",
			qs:   q,
			args: args{[]string{"a", "b"}},
			want: false,
		},
		{
			name: "Get a->b->c->d",
			qs:   q,
			args: args{[]string{"a", "b", "c", "d"}},
			want: false,
		},
		{
			name: "Get a->g",
			qs:   q,
			args: args{[]string{"a", "g"}},
			want: false,
		},
		{
			name: "Get d",
			qs:   q,
			args: args{[]string{"d"}},
			want: false,
		},
		{
			name: "Get j",
			qs:   q,
			args: args{[]string{"j"}},
			want: true,
		},
		{
			name: "Get z",
			qs:   q,
			args: args{[]string{"z"}},
			want: false,
		},
		{
			name: "Get a",
			qs:   q,
			args: args{[]string{"a"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qs.GetBool(tt.args.path...); got != tt.want {
				t.Errorf("QS.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_GetStringSlice(t *testing.T) {
	query := "a[]=b&a[]=c&a[]=d&e[f]=g&e[f]=h"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want []string
	}{
		{
			name: "Get a",
			qs:   q,
			args: args{[]string{"a"}},
			want: []string{"b", "c", "d"},
		},
		{
			name: "Get e->f",
			qs:   q,
			args: args{[]string{"e", "f"}},
			want: []string{"g", "h"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qs.GetStringSlice(tt.args.path...); !cmp.Equal(got, tt.want) {
				t.Errorf("QS.GetStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_GetAll(t *testing.T) {
	type args struct {
		path []string
	}
	tests := []struct {
		name      string
		delimiter string
		args      args
		want      []interface{}
	}{
		{
			name:      "No path",
			delimiter: "",
			args:      args{[]string{}},
			want:      []interface{}{},
		},
		{
			name:      "Get a->b",
			delimiter: "",
			args:      args{[]string{"a", "b"}},
			want:      []interface{}{"123"},
		},
		{
			name:      "Get a->g",
			delimiter: "",
			args:      args{[]string{"a", "g"}},
			want:      []interface{}{"h", "i"},
		},
		{
			name:      "Get a->g w/ delimiter",
			delimiter: ".",
			args:      args{[]string{"a.g"}},
			want:      []interface{}{"h", "i"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(
				"a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true",
				PathDelimiter(tt.delimiter),
			)
			if err != nil {
				t.Fatalf("NewQS failed with err, %s", err)
			}

			if got := q.GetAll(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_Set(t *testing.T) {
	type args struct {
		vals []interface{}
		path []string
	}
	tests := []struct {
		name      string
		delimiter string
		args      args
		want      []interface{}
	}{
		{
			name:      "Set a->b",
			delimiter: "",
			args: args{
				vals: []interface{}{"a", 123, true},
				path: []string{"a", "b"},
			},
			want: []interface{}{"a", 123, true},
		},
		{
			name:      "Set a->b w/ delimiter",
			delimiter: ".",
			args: args{
				vals: []interface{}{"a", 123, true},
				path: []string{"a.b"},
			},
			want: []interface{}{"a", 123, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(
				"a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true",
				PathDelimiter(tt.delimiter),
			)
			if err != nil {
				t.Fatalf("NewQS failed with err, %s", err)
			}

			q.Set(tt.args.vals, tt.args.path...)
			if got := q.GetAll(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_Add(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		val  interface{}
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want []interface{}
	}{
		{
			name: "Add a->b",
			qs:   q,
			args: args{
				val:  true,
				path: []string{"a", "b"},
			},
			want: []interface{}{"123", true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q.Add(tt.args.val, tt.args.path...)
			if got := q.GetAll(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseKey(t *testing.T) {
	type args struct {
		key      string
		maxDepth int
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "empty key",
			args:    args{key: ""},
			want:    []string{""},
			wantErr: false,
		},
		{
			name:    "a",
			args:    args{key: "a"},
			want:    []string{"a"},
			wantErr: false,
		},
		{
			name:    "alpha[beta][gamma]",
			args:    args{key: "alpha[beta][gamma]"},
			want:    []string{"alpha", "beta", "gamma"},
			wantErr: false,
		},
		{
			name:    "alpha[[beta]",
			args:    args{key: "alpha[[beta]"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "alpha[beta]]",
			args:    args{key: "alpha[beta]]"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "a[b][c][d] max 2",
			args:    args{key: "a[b][c][d]", maxDepth: 2},
			want:    []string{"a", "b", "[c][d]"},
			wantErr: false,
		},
		{
			name:    "a[b][c][d][e][f] max 5",
			args:    args{key: "a[b][c][d][e][f]", maxDepth: 5},
			want:    []string{"a", "b", "c", "d", "e", "[f]"},
			wantErr: false,
		},
		{
			name:    "a[b][c][d] max -1",
			args:    args{key: "a[b][c][d]", maxDepth: -1},
			want:    []string{"a", "b", "c", "d"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseKey(tt.args.key, tt.args.maxDepth)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_String(t *testing.T) {
	type fields struct {
		RawQuery string
		Values   *node
		MaxDepth int
		mutex    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1",
			fields: fields{
				mutex: &sync.RWMutex{},
				Values: &node{
					Key: "",
					Children: map[string]*node{
						"a": {
							Key:      "a",
							Values:   []interface{}{"val1"},
							Children: map[string]*node{},
						},
						"b": {
							Key:    "b",
							Values: []interface{}{"val2", "val3"},
							Children: map[string]*node{
								"c": {
									Key:      "c",
									Values:   []interface{}{"val4"},
									Children: map[string]*node{},
								},
							},
						},
					},
				},
			},
			want: "a=val1&b=val2&b=val3&b[c]=val4",
		},
		{
			name: "2",
			fields: fields{
				mutex: &sync.RWMutex{},
				Values: &node{
					Key:    "",
					Values: make([]interface{}, 0),
					Children: map[string]*node{
						"a": {
							Key:    "a",
							Values: make([]interface{}, 0),
							Children: map[string]*node{
								"b": {
									Key:    "b",
									Values: []interface{}{"123"},
									Children: map[string]*node{
										"c": {
											Key:    "c",
											Values: make([]interface{}, 0),
											Children: map[string]*node{
												"d": {
													Key:      "d",
													Values:   []interface{}{"c"},
													Children: make(map[string]*node),
												},
											},
										},
									},
								},
								"g": {
									Key:      "g",
									Values:   []interface{}{"h", "i"},
									Children: make(map[string]*node),
								},
							},
						},
						"d": {
							Key:      "d",
							Values:   []interface{}{"1.05", "2.5"},
							Children: make(map[string]*node),
						},
						"j": {
							Key:      "j",
							Values:   []interface{}{"true"},
							Children: make(map[string]*node),
						},
					},
				},
			},
			want: "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d=1.05&d=2.5&j=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QS{
				RawQuery: tt.fields.RawQuery,
				Values:   tt.fields.Values,
				MaxDepth: tt.fields.MaxDepth,
				mutex:    tt.fields.mutex,
			}

			got := q.String()
			if !assertQueryStringsEqual(got, tt.want) {
				t.Errorf("QS.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_EncodedString(t *testing.T) {
	type fields struct {
		RawQuery string
		Values   *node
		MaxDepth int
		mutex    *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1",
			fields: fields{
				mutex: &sync.RWMutex{},
				Values: &node{
					Key: "",
					Children: map[string]*node{
						"a": {
							Key:      "a",
							Values:   []interface{}{"val1"},
							Children: map[string]*node{},
						},
						"b": {
							Key:    "b",
							Values: []interface{}{"val2", "val3"},
							Children: map[string]*node{
								"c": {
									Key:      "c",
									Values:   []interface{}{"val4"},
									Children: map[string]*node{},
								},
							},
						},
					},
				},
			},
			want: "a=val1&b=val2&b=val3&b%5Bc%5D=val4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QS{
				RawQuery: tt.fields.RawQuery,
				Values:   tt.fields.Values,
				MaxDepth: tt.fields.MaxDepth,
				mutex:    tt.fields.mutex,
			}

			got := q.EncodedString()
			if !assertQueryStringsEqual(got, tt.want) {
				t.Errorf("QS.EncodedString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertQueryStringsEqual(a, b string) bool {
	aParts := strings.Split(a, "&")
	bParts := strings.Split(b, "&")

	if len(aParts) != len(bParts) {
		return false
	}

	aMap := map[string]bool{}
	for _, v := range aParts {
		aMap[v] = true
	}

	for _, v := range bParts {
		if !aMap[v] {
			return false
		}
	}

	return true
}
