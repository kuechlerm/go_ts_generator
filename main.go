package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"

	"github.com/fatih/structtag"
)

type Test struct {
	Name string `json:"name"`
}

// Converts Go type to ArkType type
func goTypeToArkType(goType string) string {
	switch goType {
	case "string":
		return `"string"`
	case "int", "int8", "int16", "int32", "int64", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64":
		return `"number"`
	case "bool":
		return `"boolean"`
	default:
		return `"any"` // fallback
	}
}

// Converts PascalCase to camelCase
func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// Generates ArkType for a struct type
func structToArkTypeSpec(typeSpec *ast.TypeSpec) string {
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return ""
	}

	var fields []string

	for _, field := range structType.Fields.List {
		fieldType := ""

		if field.Tag != nil {
			tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
			if err != nil {
				fmt.Printf("Error parsing tags for field %s: %v\n", field.Names[0].Name, err)
				continue
			}

			fmt.Printf("Processing field: %s || Tags: %s \n", field.Names, tags.Tags())
		}

		switch ft := field.Type.(type) {
		case *ast.Ident:
			fieldType = goTypeToArkType(ft.Name)

		default:
			fieldType = "any"
		}

		for _, name := range field.Names {
			// todo: json-Name statt name.Name
			fields = append(fields, fmt.Sprintf("  %s: %s", toCamelCase(name.Name), fieldType))
		}
	}

	return fmt.Sprintf("export const %s = type({\n%s\n});\n", typeSpec.Name.Name, strings.Join(fields, ",\n"))
}

// Parse Go file and return ArkType TypeScript code as string
func GenerateArkTypeFromGoStructsSrc(goSrc string) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", goSrc, parser.AllErrors)
	if err != nil {
		return "", err
	}

	var tsDefs []string
	tsDefs = append(tsDefs, `import { type } from "arktype";`)

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				tsDefs = append(tsDefs, structToArkTypeSpec(typeSpec))
			}
		}
	}

	tsCode := strings.Join(tsDefs, "\n\n")
	return tsCode, nil
}

// Reads Go source from a file
func ReadGoFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Writes ArkType TypeScript code to a file
func WriteTypeScriptFile(filename, code string) error {
	return os.WriteFile(filename, []byte(code), 0644)
}

func main() {
	// goSrc, err := ReadGoFile("input.go")
	// if err != nil {
	// 	fmt.Println("Error reading Go file:", err)
	// 	return
	// }
	// tsCode, err := GenerateArkTypeFromGoStructsSrc(goSrc)
	// if err != nil {
	// 	fmt.Println("Error generating ArkType:", err)
	// 	return
	// }
	// err = WriteTypeScriptFile("output.arktype.ts", tsCode)
	// if err != nil {
	// 	fmt.Println("Error writing TypeScript file:", err)
	// 	return
	// }
	// fmt.Println("ArkType definitions generated in output.arktype.ts")

	// -------- arktype generation example --------
	// goText := "package egal\n\n" +
	// 	"type Affe struct {\n" +
	// 	"Name string `json:\"min\"`\n" +
	// 	"Email string `json:\"email\"`\n" +
	// 	"Alter int\n" +
	// 	"}"
	//
	// tsCode, err := GenerateArkTypeFromGoStructsSrc(goText)
	// if err != nil {
	// 	fmt.Println("Error generating ArkType:", err)
	// 	return
	// }
	//
	// fmt.Println("Start Go to ArkType generator...")
	// fmt.Println(tsCode)

	text, err := os.ReadFile("beispiel/beispiel_handler.go")
	if err != nil {
		fmt.Println("Error reading Go file:", err)
		return
	}

	arkTypeCode, err := GenerateArkTypeFromGoStructsSrc(string(text))
	if err != nil {
		fmt.Println("Error generating ArkType:", err)
		return
	}

	fmt.Println(arkTypeCode)

	Init_server()
	// Init_Auth()
}
