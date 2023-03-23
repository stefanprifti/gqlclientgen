package gen

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPackageName = "main"

type Options struct {
	// PackageName is the name of the package to generate
	PackageName string
}

// GenerateTypes generates the types for the given schema and query
func GenerateTypes(ctx context.Context, schema *ast.Schema, query *ast.QueryDocument, options Options) ([]byte, error) {
	// check if context is done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// continue
		packageName := options.PackageName
		if packageName == "" {
			packageName = defaultPackageName
		}

		// generate the types from the schema
		schemaTypes, err := GenerateTypesFromSchema(packageName, schema)
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema types: %w", err)
		}

		// generate the types from the query
		queryTypes := GenerateTypesFromOperation(query)

		// combine the types
		var b bytes.Buffer
		b.Write(schemaTypes)
		b.Write(queryTypes)

		return b.Bytes(), nil
	}
}

func getOrderedTypes(schema *ast.Schema) []*ast.Definition {
	types := make([]*ast.Definition, 0, len(schema.Types))

	for _, t := range schema.Types {
		if t != nil && t.Name != "" {
			types = append(types, t)
		}
	}

	// order the types by name
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	return types
}

func GenerateTypesFromSchema(packageName string, schema *ast.Schema) ([]byte, error) {
	var b bytes.Buffer

	// write the package name
	b.WriteString("package " + packageName + "\n\n")

	for _, t := range getOrderedTypes(schema) {
		// skip the built in types
		if t.BuiltIn {
			continue
		}

		// skip the root types
		if t.Name == "Query" || t.Name == "Mutation" || t.Name == "Subscription" {
			continue
		}

		// skip the introspection types
		if strings.HasPrefix(t.Name, "__") || strings.HasPrefix(t.Name, "_") {
			continue
		}

		switch t.Kind {
		case ast.Scalar:
			// todo handle scalars
			b.WriteString("type " + t.Name + " string\n")

		case ast.Object, ast.Interface, ast.Union, ast.InputObject:
			// write the struct definition
			b.WriteString("type " + t.Name + " struct {\n")

			for _, f := range t.Fields {
				if f.Description != "" {
					b.WriteString("\t// " + f.Description + "\n")
				}

				// write the field name and type
				printFieldDefinition(f, &b, 1)
			}

			b.WriteString("}\n")
		case ast.Enum:
			b.WriteString("type " + t.Name + " string\n")
			b.WriteString("const (\n")
			for _, v := range t.EnumValues {
				b.WriteString("\t" + v.Name + " " + t.Name + " = \"" + v.Name + "\"\n")
			}
			b.WriteString(")\n")

		default:
			continue
		}
	}

	return b.Bytes(), nil
}

func GenerateTypesFromOperation(doc *ast.QueryDocument) []byte {
	var b bytes.Buffer

	for _, op := range doc.Operations {
		// Print the request struct
		fmt.Fprintf(&b, "type %sRequest struct {\n", op.Name)
		for _, v := range op.VariableDefinitions {
			if v.Type.NonNull {
				fmt.Fprintf(&b, "\t%s %s `json:\"%s\"`\n", toCammelCase(v.Variable), convertGraphQLTypeToGoType(v.Type), v.Variable)
			} else {
				fmt.Fprintf(&b, "\t%s %s `json:\"%s,omitempty\"`\n", toCammelCase(v.Variable), convertGraphQLTypeToGoType(v.Type), v.Variable)
			}
		}
		fmt.Fprintln(&b, "}")

		// Print the response struct
		fmt.Fprintf(&b, "type %sResponse struct {\n", op.Name)

		// TODO: maybe add option to skip the first selection set?
		generateResponseTypes(op.SelectionSet, &b, 1)
		fmt.Fprintln(&b, "}")
	}
	return b.Bytes()
}

func generateResponseTypes(sel ast.SelectionSet, b *bytes.Buffer, level int) {
	for _, s := range sel {
		switch s := s.(type) {
		case *ast.Field:
			if isPrimitiveGQLType(s.Definition.Type) {
				printFieldDefinition(s.Definition, b, level)
			} else {
				fmt.Fprintf(b, "%s%s struct {\n", strings.Repeat("\t", level), toCammelCase(s.Alias))
				if s.SelectionSet != nil {
					generateResponseTypes(s.SelectionSet, b, level+1)
				}
				if s.Definition.Type.NonNull {
					fmt.Fprintf(b, "%s} `json:\"%s\"`\n", strings.Repeat("\t", level), s.Alias)
				} else {
					fmt.Fprintf(b, "%s} `json:\"%s,omitempty\"`\n", strings.Repeat("\t", level), s.Alias)
				}
			}
		case *ast.InlineFragment:
			panic("inline fragment not supported")
		case *ast.FragmentSpread:
			panic("fragment spread not supported")
		}
	}
}

func printFieldDefinition(f *ast.FieldDefinition, b *bytes.Buffer, level int) {
	if f.Type.NonNull {
		// todo: add alias
		fmt.Fprintf(b, "%s%s %s `json:\"%s\"`\n", strings.Repeat("\t", level), toCammelCase(f.Name), convertGraphQLTypeToGoType(f.Type), f.Name)
	} else {
		fmt.Fprintf(b, "%s%s %s `json:\"%s,omitempty\"`\n", strings.Repeat("\t", level), toCammelCase(f.Name), convertGraphQLTypeToGoType(f.Type), f.Name)
	}
}

func convertGraphQLTypeToGoType(typ *ast.Type) string {
	if typ.Elem != nil {
		return "[]" + convertGraphQLTypeToGoType(typ.Elem)
	}
	switch typ.Name() {
	case "String":
		return "string"
	case "Int":
		return "int"
	case "Float":
		return "float64"
	case "Boolean":
		return "bool"
	case "ID":
		return "string"
	default:
		return typ.Name()
	}
}

func isPrimitiveGQLType(typ *ast.Type) bool {
	switch typ.Name() {
	case "String", "Int", "Float", "Boolean", "ID":
		return true
	default:
		return false
	}
}

func toCammelCase(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
