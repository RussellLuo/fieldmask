package fieldmask_test

import (
	"reflect"
	"testing"

	"github.com/RussellLuo/fieldmask"
)

var testMap = map[string]interface{}{
	"name": "foo",
	"age":  1,
	"address": map[string]interface{}{
		"country": "X",
		"city":    "Y",
	},
}

func TestGet(t *testing.T) {
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
			wantValue: 1,
			wantOK:    true,
		},
		{
			name:      "is_male absent",
			inPath:    "is_male",
			wantValue: nil,
			wantOK:    false,
		},
		{
			name:   "address exists",
			inPath: "address",
			wantValue: map[string]interface{}{
				"country": "X",
				"city":    "Y",
			},
			wantOK: true,
		},
		{
			name:      "address.country exists",
			inPath:    "address.country",
			wantValue: "X",
			wantOK:    true,
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
			wantValue: "Y",
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

func TestCopy(t *testing.T) {
	src := testMap

	tests := []struct {
		name    string
		inPaths []string
		wantDst map[string]interface{}
	}{
		{
			name:    "copy name",
			inPaths: []string{"name"},
			wantDst: map[string]interface{}{
				"name": "foo",
			},
		},
		{
			name:    "copy age",
			inPaths: []string{"age"},
			wantDst: map[string]interface{}{
				"age": 1,
			},
		},
		{
			name:    "copy is_male",
			inPaths: []string{"is_male"},
			wantDst: map[string]interface{}{},
		},
		{
			name:    "copy address",
			inPaths: []string{"address"},
			wantDst: map[string]interface{}{
				"address": map[string]interface{}{
					"country": "X",
					"city":    "Y",
				},
			},
		},
		{
			name:    "copy address.country",
			inPaths: []string{"address.country"},
			wantDst: map[string]interface{}{
				"address": map[string]interface{}{
					"country": "X",
				},
			},
		},
		{
			name:    "copy address.province",
			inPaths: []string{"address.province"},
			wantDst: map[string]interface{}{
				"address": map[string]interface{}{},
			},
		},
		{
			name:    "copy address.city",
			inPaths: []string{"address.city"},
			wantDst: map[string]interface{}{
				"address": map[string]interface{}{
					"city": "Y",
				},
			},
		},
		{
			name:    "copy multiple fields",
			inPaths: []string{"name", "age", "is_male", "address.city"},
			wantDst: map[string]interface{}{
				"name": "foo",
				"age":  1,
				"address": map[string]interface{}{
					"city": "Y",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm := fieldmask.FieldMask{}
			fm.Copy(src, tt.inPaths...)

			gotDst := map[string]interface{}(fm)

			if !reflect.DeepEqual(gotDst, tt.wantDst) {
				t.Fatalf("Dst: Got (%#v) != Want (%#v)", gotDst, tt.wantDst)
			}
		})
	}
}
