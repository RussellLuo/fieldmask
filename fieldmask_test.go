package fieldmask_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/fieldmask"
)

var testMap = map[string]interface{}{
	"name": "foo",
	"age":  20,
	"address": map[string]interface{}{
		"city": "Z",
	},
}

func TestFrom(t *testing.T) {
	src := testMap

	tests := []struct {
		name    string
		inPaths []string
		want    fieldmask.FieldMask
	}{
		{
			name:    "select name",
			inPaths: []string{"name"},
			want: fieldmask.FieldMask{
				"name": "foo",
			},
		},
		{
			name:    "select age",
			inPaths: []string{"age"},
			want: fieldmask.FieldMask{
				"age": 20,
			},
		},
		{
			name:    "select xxx",
			inPaths: []string{"xxx"},
			want: fieldmask.FieldMask{
				"xxx": nil,
			},
		},
		{
			name:    "select address",
			inPaths: []string{"address"},
			want: fieldmask.FieldMask{
				"address": map[string]interface{}{
					"city": "Z",
				},
			},
		},
		{
			name:    "select address.country",
			inPaths: []string{"address.country"},
			want: fieldmask.FieldMask{
				"address": map[string]interface{}{
					"country": nil,
				},
			},
		},
		{
			name:    "select address.province",
			inPaths: []string{"address.province"},
			want: fieldmask.FieldMask{
				"address": map[string]interface{}{
					"province": nil,
				},
			},
		},
		{
			name:    "select address.city",
			inPaths: []string{"address.city"},
			want: fieldmask.FieldMask{
				"address": map[string]interface{}{
					"city": "Z",
				},
			},
		},
		{
			name:    "select deep1.deep2.deep3.key",
			inPaths: []string{"deep1.deep2.deep3.key"},
			want: fieldmask.FieldMask{
				"deep1": map[string]interface{}{
					"deep2": map[string]interface{}{
						"deep3": map[string]interface{}{
							"key": nil,
						},
					},
				},
			},
		},
		{
			name:    "select multiple fields",
			inPaths: []string{"name", "xxx", "address.city"},
			want: fieldmask.FieldMask{
				"name": "foo",
				"xxx":  nil,
				"address": map[string]interface{}{
					"city": "Z",
				},
			},
		},
		{
			name:    "select all fields",
			inPaths: nil,
			want: fieldmask.FieldMask{
				"name": "foo",
				"age":  20,
				"address": map[string]interface{}{
					"city": "Z",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldmask.From(src, tt.inPaths...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("FieldMask: Got (%#v) != Want (%#v)", got, tt.want)
			}
		})
	}
}

func TestFieldMask_Get(t *testing.T) {
	fm := fieldmask.FieldMask(testMap)

	tests := []struct {
		name      string
		inPath    string
		wantValue interface{}
		wantOK    bool
	}{
		{
			name:      "name exists",
			inPath:    "name",
			wantValue: "foo",
			wantOK:    true,
		},
		{
			name:      "age exists",
			inPath:    "age",
			wantValue: 20,
			wantOK:    true,
		},
		{
			name:   "address exists",
			inPath: "address",
			wantValue: map[string]interface{}{
				"city": "Z",
			},
			wantOK: true,
		},
		{
			name:      "address.country absent",
			inPath:    "address.country",
			wantValue: nil,
			wantOK:    false,
		},
		{
			name:      "address.province absent",
			inPath:    "address.province",
			wantValue: nil,
			wantOK:    false,
		},
		{
			name:      "address.city exists",
			inPath:    "address.city",
			wantValue: "Z",
			wantOK:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOK := fm.Get(tt.inPath)

			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Fatalf("Value: Got (%#v) != Want (%#v)", gotValue, tt.wantValue)
			}

			if gotOK != tt.wantOK {
				t.Fatalf("OK: Got (%#v) != Want (%#v)", gotOK, tt.wantOK)
			}
		})
	}
}
