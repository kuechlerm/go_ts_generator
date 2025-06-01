package clientGenerator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"testing"
)

func TestListStructs(t *testing.T) {
	file_content, err := os.ReadFile("../beispiel/beispiel_handler.go")
	if err != nil {
		t.Fatalf("Error reading Go file: %v", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file_content, parser.AllErrors)
	if err != nil {
		t.Fatalf("Error parsing Go file: %v", err)
	}

	var structs []string
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				// tsDefs = append(tsDefs, structToArkTypeSpec(typeSpec))
				structs = append(structs, typeSpec.Name.Name)
			}
		}
	}

	expected := []string{
		"affe",
	}

	if !reflect.DeepEqual(structs, expected) {
		t.Errorf("slices are not equal: %v vs %v", structs, expected)
	}
}
