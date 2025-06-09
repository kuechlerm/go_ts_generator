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

func Run(go_file_path, target_path string) error {
	go_content, err := os.ReadFile(go_file_path)
	if err != nil {
		return errors.New("Error reading Go file: " + err.Error())
	}

	ts_code, err := generateTS(string(go_content))
	if err != nil {
		return errors.New("Error generating TypeScript code: " + err.Error())
	}

	err = os.WriteFile(target_path, []byte(ts_code), 0644)
	if err != nil {
		return errors.New("Error writing TypeScript file: " + err.Error())
	}

	return nil
}

func getRPCs(file_content string) ([]RPC, error) {
	rpcs := []RPC{}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", file_content, parser.AllErrors)
	if err != nil {
		return rpcs, errors.New("Error parsing Go file: " + err.Error())
	}

	rpcNameMap := map[string]RPC{}

	for _, decl := range node.Decls {
		// fmt.Println(decl)
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
					call.request = mapSchema(typeSpec)
				}

				if strings.HasSuffix(typeSpec.Name.Name, "_Response") {
					call.response = mapSchema(typeSpec)
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

func generateTS(input string) (string, error) {
	rpcs, err := getRPCs(input)
	if err != nil {
		return "", err
	}

	var tsCode strings.Builder
	tsCode.WriteString(`import { type } from "arktype";`)
	tsCode.WriteString("\n\n")

	for _, rpc := range rpcs {
		// generate request schema
		tsCode.WriteString(fmt.Sprintf("export const %s_Schema = type({\n", rpc.request.Name))
		for _, prop := range rpc.request.Properties {
			tsCode.WriteString(fmt.Sprintf(`  %s: "%s",`, prop.Name, prop.Type))
			tsCode.WriteString("\n")
		}
		tsCode.WriteString("});\n\n")

		// generate request type
		tsCode.WriteString(fmt.Sprintf("export type %s = typeof %s_Schema.infer;\n\n", rpc.request.Name, rpc.request.Name))

		// generate response schema
		tsCode.WriteString(fmt.Sprintf("export const %s_Schema = type({\n", rpc.response.Name))
		for _, prop := range rpc.response.Properties {
			tsCode.WriteString(fmt.Sprintf(`  %s: "%s",`, prop.Name, prop.Type))
			tsCode.WriteString("\n")
		}
		tsCode.WriteString("});\n\n")

		// generate response type
		tsCode.WriteString(fmt.Sprintf("export type %s = typeof %s_Schema.infer;\n\n", rpc.response.Name, rpc.response.Name))
	}

	// rpc client class
	tsCode.WriteString("export class RPC_Client {\n")
	tsCode.WriteString("  constructor(public base_url: string) {}\n\n")
	tsCode.WriteString("  async #do_fetch<TRequest, TResponse>(\n")
	tsCode.WriteString("    path: string,\n")
	tsCode.WriteString("    args: TRequest,\n")
	tsCode.WriteString("  ): Promise<{ result: TResponse | null; error: string | null }> {\n")
	tsCode.WriteString("    try {\n")
	tsCode.WriteString("      const result = await fetch(new URL(path, this.base_url).href, {\n")
	tsCode.WriteString("        method: \"POST\",\n")
	tsCode.WriteString("        body: JSON.stringify(args),\n")
	tsCode.WriteString("      });\n\n")
	tsCode.WriteString("      if (!result.ok) {\n")
	tsCode.WriteString("        console.error(\n")
	tsCode.WriteString("          `Fetch error: ${result.status} ${result.statusText} for ${path}`,\n")
	tsCode.WriteString("        );\n")
	tsCode.WriteString("        return {\n")
	tsCode.WriteString("          result: null,\n")
	tsCode.WriteString("          error: `Fetch error: ${result.status} ${result.statusText}`,\n")
	tsCode.WriteString("        };\n")
	tsCode.WriteString("      }\n\n")
	tsCode.WriteString("      const data = await result.json();\n\n")
	tsCode.WriteString("      return {\n")
	tsCode.WriteString("        result: data as TResponse,\n")
	tsCode.WriteString("        error: null,\n")
	tsCode.WriteString("      };\n")
	tsCode.WriteString("    } catch (error) {\n")
	tsCode.WriteString("      console.error(`Error during fetch for ${path}:`, error);\n\n")
	tsCode.WriteString("      return {\n")
	tsCode.WriteString("        result: null,\n")
	tsCode.WriteString("        error: error instanceof Error ? error.message : \"Unknown error\",\n")
	tsCode.WriteString("      };\n")
	tsCode.WriteString("    }\n")
	tsCode.WriteString("  }\n\n")

	for _, rpc := range rpcs {
		tsCode.WriteString(
			"  " +
				strings.ToLower(strings.Split(rpc.request.Name, "_")[0]) +
				" = (args: " + rpc.request.Name + ") =>\n")
		tsCode.WriteString(
			"    this.#do_fetch<" +
				rpc.request.Name +
				", " +
				rpc.response.Name +
				">(\"" + rpc.path + "\", args);\n")
	}

	tsCode.WriteString("}\n")

	return tsCode.String(), nil
}

func mapSchema(typeSpec *ast.TypeSpec) Schema {
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

			// fmt.Printf("Processing field: %s || Tags: %s \n", field.Names, tags.Tags())

			validate_tag_used := false
			for _, tag := range tags.Tags() {
				if tag.Key == "json" {
					jsonPropertyName = tag.Name // ist der erste Tag-Wert
				}

				if tag.Key == "validate" {
					// fmt.Printf("Validation tag found: %s\n", tag.Name)
					// fmt.Printf("Validation OPTIONS tag found: %s\n", tag.Options)
					fieldType = mapValidation(fieldType, tag.Name)
					validate_tag_used = true
				}
			}

			if !validate_tag_used {
				fieldType = mapValidation(fieldType, "")
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

func mapValidation(ts_typ, validation string) string {
	if ts_typ == "string" || ts_typ == "number" {
		switch validation {
		case "required":
			return ts_typ + " > 0"
		case "":
			return ts_typ + " | undefined"
		}
	}

	if ts_typ == "boolean" {
		switch validation {
		case "required":
			return "true"
		case "":
			return ts_typ + " | undefined"
		}
	}

	return "todo"
}
