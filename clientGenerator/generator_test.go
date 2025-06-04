package clientGenerator

import (
	"reflect"
	"testing"
)

func TestGetRPCInfos(t *testing.T) {
	rpcs, err := GetRPCs("../beispiel/beispiel_handler.go")
	if err != nil {
		t.Fatalf("Error getting RPCs: %v", err)
	}

	expected := []RPC{
		{
			path: "/beispielanlegen",
			request: Schema{
				Name: "BeispielAnlegen_Request",
				Properties: []Property{
					{"name", "string", "3 <= string <= 100"},
				},
			},
			response: Schema{
				Name: "BeispielAnlegen_Response",
				Properties: []Property{
					{"id", "string", ""},
				},
			},
		},
		{
			path: "/beispielaendern",
			request: Schema{
				Name: "BeispielAendern_Request",
				Properties: []Property{
					{"name", "string", "3 <= string <= 100"},
				},
			},
			response: Schema{
				Name: "BeispielAendern_Response",
				Properties: []Property{
					{"id", "string", ""},
				},
			},
		},
	}

	if !reflect.DeepEqual(rpcs, expected) {
		t.Errorf("slices are not equal: %v vs %v", rpcs, expected)
	}
}

func TestMapValidation(t *testing.T) {
	tests := []struct {
		typ      string
		validate string
		arktype  string
	}{
		{"string", "required", "string > 0"},
		{"string", "", "string | undefined"},

		{"number", "required", "number > 0"},
		{"number", "", "number | undefined"},

		{"boolean", "required", "boolean"},
		{"boolean", "", "boolean | undefined"},

		// todo: date, arrays, ... -> Amplenote
	}

	for _, test := range tests {
		result := MapValidation(test.typ, test.validate)
		if result != test.arktype {
			t.Errorf("MapValidation(%q, %q) = %q; want %q", test.typ, test.validate, result, test.arktype)
		}
	}
}
