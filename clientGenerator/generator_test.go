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
					{Name: "name", Type: "string", Validation: "3 <= string <= 100"},
				},
			},
			response: Schema{
				Name: "BeispielAnlegen_Response",
				Properties: []Property{
					{Name: "id", Type: "string", Validation: ""},
				},
			},
		},
		{
			path: "/beispielaendern",
			request: Schema{
				Name: "BeispielAendern_Request",
				Properties: []Property{
					{Name: "name", Type: "string", Validation: "3 <= string <= 100"},
				},
			},
			response: Schema{
				Name: "BeispielAendern_Response",
				Properties: []Property{
					{Name: "id", Type: "string", Validation: ""},
				},
			},
		},
	}

	if !reflect.DeepEqual(rpcs, expected) {
		t.Errorf("slices are not equal: %v vs %v", rpcs, expected)
	}
}

func TestMapValidation(t *testing.T) {
	// tests := []struct {
	// 	input    string
	// 	expected string
	// }{
	// 	{"required", "required"},
	// 	{"min=3", "min:3"},
	// 	{"max=100", "max:100"},
	// 	{"required,min=3,max=100", "required,min:3,max:100"},
	// }
	//
	// for _, test := range tests {
	// 	result := MapValidation(test.input)
	// 	if result != test.expected {
	// 		t.Errorf("MapValidation(%q) = %q; want %q", test.input, result, test.expected)
	// 	}
	// }
}
