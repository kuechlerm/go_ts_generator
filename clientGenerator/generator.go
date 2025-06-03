package clientGenerator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
)

type Property struct {
	Name       string // json Name
	Type       string // TS Type
	Validation string // Ark Validation
}

type Schema struct {
	Name       string
	Properties []Property
}

type RPC struct {
	path     string
	request  Schema
	response Schema
}

func GetRPCs(file_path string) ([]RPC, error) {
	rpcs := []RPC{}

	file_content, err := os.ReadFile("../beispiel/beispiel_handler.go")
	if err != nil {
		return rpcs, errors.New("Error reading Go file: " + err.Error())
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file_content, parser.AllErrors)
	if err != nil {
		return rpcs, errors.New("Error parsing Go file: " + err.Error())
	}

	rpcNameMap := map[string]RPC{}

	for _, decl := range node.Decls {
		fmt.Println(decl)
		genDecl, ok := decl.(*ast.GenDecl)

		if !ok || (genDecl.Tok != token.TYPE && genDecl.Tok != token.CONST) {
			// wir interessieren uns nur für Typ- und Konstantendeklarationen
			continue
		}

		for _, spec := range genDecl.Specs {
			constSpec, ok := spec.(*ast.ValueSpec)

			if ok {
				constName := constSpec.Names[0].Name
				// muss gleich specName sein, damit Zuordnung stimmt
				constSpacName := strings.Split(constName, "_")[0]

				// Wir suchen nach Konstanten, die mit "_Path" enden
				literal, ok := constSpec.Values[0].(*ast.BasicLit)
				// todo: check für Names[0]?
				if !ok || literal.Kind != token.STRING || !strings.HasSuffix(constName, "_Path") {
					continue
				}

				rpc := rpcNameMap[constSpacName]
				rpc.path = strings.Trim(literal.Value, "\"")
				rpcNameMap[constSpacName] = rpc

				continue
			}

			typeSpec, ok := spec.(*ast.TypeSpec)
			specName := strings.Split(typeSpec.Name.Name, "_")[0]
			if !ok {
				continue
			}

			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				call := rpcNameMap[specName]

				if strings.HasSuffix(typeSpec.Name.Name, "_Request") {
					call.request = MapSchema(typeSpec)
				}

				if strings.HasSuffix(typeSpec.Name.Name, "_Response") {
					call.response = MapSchema(typeSpec)
				}

				rpcNameMap[specName] = call
			}
		}
	}

	for _, call := range rpcNameMap {
		rpcs = append(rpcs, call)
	}

	return rpcs, nil
}

func MapSchema(typeSpec *ast.TypeSpec) Schema {
	properties := []Property{}

	for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {

		// ##### Type
		fieldType := ""
		switch ft := field.Type.(type) {
		case *ast.Ident:
			fieldType = goTypeToArkType(ft.Name)
		default:
			fieldType = "any"
		}

		// ##### Tags
		jsonPropertyName := ""
		if field.Tag != nil {

			tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
			if err != nil {
				fmt.Printf("Error parsing tags for field %s: %v\n", field.Names[0].Name, err)
				continue
			}

			fmt.Printf("Processing field: %s || Tags: %s \n", field.Names, tags.Tags())

			for _, tag := range tags.Tags() {
				if tag.Key == "json" {
					jsonPropertyName = tag.Name // ist der erste Tag-Wert
				}

				if tag.Key == "validate" {
					// todo:
					fmt.Printf("Validation tag found: %s\n", tag.Name)
					fmt.Printf("Validation OPTIONS tag found: %s\n", tag.Options)
				}
			}
		}

		name := field.Names[0].Name
		if jsonPropertyName != "" {
			name = jsonPropertyName // wenn json-Name vorhanden, dann diesen verwenden
		}
		properties = append(properties, Property{
			Name:       name, // json name
			Type:       fieldType,
			Validation: "TODO", // TODO: hier müsste die Validation aus den Struct-Tags geholt werden
		})

	}

	return Schema{
		Name:       typeSpec.Name.Name,
		Properties: properties,
	}
}

// Converts Go type to ArkType type
func goTypeToArkType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int8", "int16", "int32", "int64", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "any" // fallback
	}
}

func MapValidation(typ, validation string) string {
	return "todo"
}
