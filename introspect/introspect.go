package introspect

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	client "github.com/stefanprifti/gqlclient"
)

// printObject prints the object type.
func convertObject(t *Type) string {
	var sb strings.Builder

	// type Name {
	typeDeclaration(&sb, t)
	openingBrace(&sb)

	for _, f := range t.Fields {
		// # Description
		fieldDescription(&sb, f.Description)

		// fieldName(argName: ArgType): FieldType
		fieldNameType(&sb, f)

		// @deprecated(reason: "Use something else.")
		deprecatedDirective(&sb, f.IsDeprecated, f.DeprecationReason)

		// \n
		newLine(&sb)
	}

	// }
	closingBrace(&sb)

	return sb.String()
}

func convertInputObject(t *Type) string {
	var sb strings.Builder

	// # Description
	typeDescription(&sb, t)

	// input Name {
	inputDeclaration(&sb, t)
	openingBrace(&sb)

	for _, f := range t.InputFields {
		// # Description
		fieldDescription(&sb, f.Description)

		// fieldName: FieldType
		inputField(&sb, f)

		// \n
		newLine(&sb)
	}

	// }
	closingBrace(&sb)

	return sb.String()
}

func convertEnum(t *Type) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("enum %s {\n", *t.Name))

	for _, f := range t.EnumValues {
		if f.Description != nil {
			sb.WriteString(fmt.Sprintf("\t# %s\n", *f.Description))
		}
		if f.IsDeprecated {
			sb.WriteString(fmt.Sprintf("\t# Deprecated: %s\n", *f.DeprecationReason))
		}
		sb.WriteString(fmt.Sprintf("\t%s\n", f.Name))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func convertScalar(t *Type) string {
	if isPrimitiveScalar(t) {
		return ""
	}

	var sb strings.Builder

	if t.Description != nil {
		sb.WriteString(fmt.Sprintf("# %s\n", *t.Description))
	}

	sb.WriteString(fmt.Sprintf("scalar %s\n", *t.Name))

	return sb.String()
}

func convertInterface(t *Type) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("interface %s {\n", *t.Name))

	// TODO: add implemented interfaces

	for _, f := range t.Fields {
		if f.Description != nil {
			sb.WriteString(fmt.Sprintf("# %s\n", *f.Description))
		}

		if f.IsDeprecated {
			sb.WriteString(fmt.Sprintf("# Deprecated: %s\n", *f.DeprecationReason))
		}

		sb.WriteString(fmt.Sprintf("\t%s: %s\n", f.Name, f.Type.String()))

	}

	sb.WriteString("}\n")

	return sb.String()

}

func convertUnion(t *Type) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("union %s = ", *t.Name))

	for i, m := range t.PossibleTypes {
		sb.WriteString(*m.Name)
		if i < len(t.PossibleTypes)-1 {
			sb.WriteString(" | ")
		}
	}

	sb.WriteString("\n")

	return sb.String()
}

func deprecatedDirective(sb *strings.Builder, isDeprecated bool, deprecationReason *string) {
	if !isDeprecated {
		return
	}

	sb.WriteString(" @deprecated")

	if deprecationReason != nil && *deprecationReason != "" {
		sb.WriteString(fmt.Sprintf("(reason: \"%s\")", *deprecationReason))
	}
}

func newLine(sb *strings.Builder) {
	sb.WriteString("\n")
}

func closingBrace(sb *strings.Builder) {
	sb.WriteString("}\n")
}

func openingBrace(sb *strings.Builder) {
	sb.WriteString("{\n")
}

func typeDeclaration(sb *strings.Builder, t *Type) {
	sb.WriteString(fmt.Sprintf("type %s ", *t.Name))
}

func fieldNameType(sb *strings.Builder, f Field) {
	var args strings.Builder

	if len(f.Args) > 0 {
		args.WriteString("(")
		for idx, a := range f.Args {
			if idx > 0 {
				args.WriteString(", ")
			}
			args.WriteString(fmt.Sprintf("%s: %s", a.Name, a.Type.String()))
			if a.DefaultValue != nil {
				args.WriteString(fmt.Sprintf(" = %s", *a.DefaultValue))
			}
		}
		args.WriteString(")")
	}

	sb.WriteString(fmt.Sprintf("\t%s%s: %s", f.Name, args.String(), f.Type.String()))
}

func fieldDescription(sb *strings.Builder, description *string) {
	if description != nil && *description != "" {
		for _, line := range strings.Split(*description, "\n") {
			sb.WriteString(fmt.Sprintf("\t# %s", line))
			newLine(sb)
		}
	}
}

func inputDeclaration(sb *strings.Builder, t *Type) {
	sb.WriteString(fmt.Sprintf("input %s", *t.Name))
}

func typeDescription(sb *strings.Builder, t *Type) {
	if t.Description != nil && *t.Description != "" {
		for _, line := range strings.Split(*t.Description, "\n") {
			sb.WriteString(fmt.Sprintf("# %s", line))
			newLine(sb)
		}
	}
}

func inputField(sb *strings.Builder, f InputField) {
	sb.WriteString(fmt.Sprintf("\t%s: %s", f.Name, f.Type.String()))

	if f.DefaultValue != nil {
		sb.WriteString(fmt.Sprintf(" = %s", *f.DefaultValue))
	}
}

// isPrimitiveScalar returns true if the type is a primitive scalar.
func isPrimitiveScalar(t *Type) bool {
	if t == nil {
		return false
	}

	if t.Kind == nil {
		return false
	}

	switch *t.Kind {
	case TypeKindScalar:
		return *t.Name == "String" || *t.Name == "Int" || *t.Name == "Float" || *t.Name == "Boolean" || *t.Name == "ID"
	default:
		return false
	}
}

// writeSchema writes the schema to the writer.
func writeSchema(schema *Schema, w io.Writer) (err error) {
	// recover from panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not generate schema: %v", r)
		}
	}()

	for _, t := range schema.Types {
		switch *t.Kind {
		case TypeKindScalar:
			_, err := w.Write([]byte(convertScalar(&t)))
			if err != nil {
				return fmt.Errorf("failed to write scalar: %w", err)
			}
		case TypeKindObject:
			_, err := w.Write([]byte(convertObject(&t)))
			if err != nil {
				return fmt.Errorf("failed to write object: %w", err)
			}

		case TypeKindInputObject:
			_, err := w.Write([]byte(convertInputObject(&t)))
			if err != nil {
				return fmt.Errorf("failed to write input object: %w", err)
			}

		case TypeKindEnum:
			_, err := w.Write([]byte(convertEnum(&t)))
			if err != nil {
				return fmt.Errorf("failed to write enum: %w", err)
			}

		case TypeKindInterface:
			_, err := w.Write([]byte(convertInterface(&t)))
			if err != nil {
				return fmt.Errorf("failed to write interface: %w", err)
			}

		case TypeKindUnion:
			_, err := w.Write([]byte(convertUnion(&t)))
			if err != nil {
				return fmt.Errorf("failed to write union: %w", err)
			}

		default:
			panic(fmt.Sprintf("unknown type kind: %s", *t.Kind))
		}
	}

	return nil
}

// URL returns the schema from the given URL.
func URL(url string) (*Schema, error) {
	var schema struct {
		Schema `json:"__schema"`
	}

	gqlClient := client.New(client.Options{
		Endpoint: url,
	})
	err := gqlClient.Query(context.Background(), introspectionQuery, nil, &schema)
	if err != nil {
		return nil, err
	}

	return &schema.Schema, nil
}

func SchemaToText(schema *Schema) ([]byte, error) {
	var buf bytes.Buffer
	err := writeSchema(schema, &buf)
	if err != nil {
		return []byte(nil), err
	}

	return buf.Bytes(), nil
}
