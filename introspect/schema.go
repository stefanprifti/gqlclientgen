package introspect

import "fmt"

type TypeKind string

const (
	TypeKindScalar      TypeKind = "SCALAR"
	TypeKindObject      TypeKind = "OBJECT"
	TypeKindInterface   TypeKind = "INTERFACE"
	TypeKindUnion       TypeKind = "UNION"
	TypeKindEnum        TypeKind = "ENUM"
	TypeKindInputObject TypeKind = "INPUT_OBJECT"
	TypeKindList        TypeKind = "LIST"
	TypeKindNonNull     TypeKind = "NON_NULL"
)

// Field represents a GraphQL field.
type Field struct {
	Name              string       `json:"name"`
	Description       *string      `json:"description"`
	Args              []InputValue `json:"args"`
	Type              Type         `json:"type"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason *string      `json:"deprecationReason"`
}

// InputValue represents a GraphQL input value.
type InputValue struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Type         Type    `json:"type"`
	DefaultValue *string `json:"defaultValue"`
}

// EnumValue represents a GraphQL enum value.
type EnumValue struct {
	Name              string  `json:"name"`
	Description       *string `json:"description"`
	IsDeprecated      bool    `json:"isDeprecated"`
	DeprecationReason *string `json:"deprecationReason"`
}

// InputField represents a GraphQL input field.
type InputField struct {
	Name         string  `json:"name"`
	Description  *string `json:"description"`
	Type         Type    `json:"type"`
	DefaultValue *string `json:"defaultValue"`
}

// Type represents a GraphQL type.
type Type struct {
	Kind          *TypeKind    `json:"kind"`
	Name          *string      `json:"name"`
	Description   *string      `json:"description"`
	Fields        []Field      `json:"fields"`
	InputFields   []InputField `json:"inputFields"`
	Interfaces    []Type       `json:"interfaces"`
	EnumValues    []EnumValue  `json:"enumValues"`
	PossibleTypes []Type       `json:"possibleTypes"`
	OfType        *Type        `json:"ofType,omitempty"`
}

// Directive represents a GraphQL directive.
type Directive struct {
	Name        string       `json:"name"`
	Description *string      `json:"description"`
	Locations   []string     `json:"locations"`
	Args        []InputValue `json:"args"`
}

// Schema represents the GraphQL schema.
type Schema struct {
	QueryType        *Type       `json:"queryType,omitempty"`
	MutationType     *Type       `json:"mutationType,omitempty"`
	SubscriptionType *Type       `json:"subscriptionType,omitempty"`
	Types            []Type      `json:"types"`
	Directives       []Directive `json:"directives"`
}

func (t *Type) String() string {
	if t == nil {
		return ""
	}

	if t.Kind == nil {
		return ""
	}

	switch *t.Kind {
	case TypeKindNonNull:
		return fmt.Sprintf("%s!", t.OfType)
	case TypeKindList:
		return fmt.Sprintf("[%s]", t.OfType)
	default:
		return *t.Name
	}
}
