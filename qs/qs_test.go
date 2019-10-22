package qs

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewQS(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	exp := &node{
		Key:    "",
		Values: make([]interface{}, 0),
		Children: map[string]*node{
			"a": &node{
				Key:    "a",
				Values: make([]interface{}, 0),
				Children: map[string]*node{
					"b": &node{
						Key:    "b",
						Values: []interface{}{"123"},
						Children: map[string]*node{
							"c": &node{
								Key:    "c",
								Values: make([]interface{}, 0),
								Children: map[string]*node{
									"d": &node{
										Key:      "d",
										Values:   []interface{}{"c"},
										Children: make(map[string]*node),
									},
								},
							},
						},
					},
					"g": &node{
						Key:      "g",
						Values:   []interface{}{"h", "i"},
						Children: make(map[string]*node),
					},
				},
			},
			"d": &node{
				Key:      "d",
				Values:   []interface{}{"1.05"},
				Children: make(map[string]*node),
			},
			"j": &node{
				Key:      "j",
				Values:   []interface{}{"true"},
				Children: make(map[string]*node),
			},
		}}

	if !(cmp.Equal(q.Values, exp)) {
		t.Fatalf("Generated incorrect tree: %s", cmp.Diff(q.Values, exp))
	}
}

func TestQS_Get(t *testing.T) {
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
		want interface{}
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
			want: nil,
		},
		{
			name: "Get a",
			qs:   q,
			args: args{[]string{"a"}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qs.Get(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
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

func TestQS_GetAll(t *testing.T) {
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
		want []interface{}
	}{
		{
			name: "Get a->b",
			qs:   q,
			args: args{[]string{"a", "b"}},
			want: []interface{}{"123"},
		},
		{
			name: "Get a->g",
			qs:   q,
			args: args{[]string{"a", "g"}},
			want: []interface{}{"h", "i"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := q.GetAll(tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QS.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQS_Set(t *testing.T) {
	query := "a[b]=123&a[b][c][d]=c&a[g]=h&a[g]=i&d[]=1.05&j=true"
	q, err := New(query)
	if err != nil {
		t.Fatalf("NewQS failed with err, %s", err)
	}

	type args struct {
		vals []interface{}
		path []string
	}
	tests := []struct {
		name string
		qs   *QS
		args args
		want []interface{}
	}{
		{
			name: "Set a->b",
			qs:   q,
			args: args{
				vals: []interface{}{"a", 123, true},
				path: []string{"a", "b"},
			},
			want: []interface{}{"a", 123, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
