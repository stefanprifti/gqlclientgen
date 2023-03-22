package gen_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stefanprifti/gql/gen"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

func loadFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestGen(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		schema   string
		expected string
	}{
		{
			name:     "simple query",
			query:    "./testdata/query.graphql",
			schema:   "./testdata/schema.graphql",
			expected: "./testdata/types.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load data
			queryContent, err := loadFile(tt.query)
			if err != nil {
				t.Errorf("could not load query file: %v", err)
				return
			}

			schemaContent, err := loadFile(tt.schema)
			if err != nil {
				t.Errorf("could not load schema file: %v", err)
				return
			}

			expectedTypes, err := loadFile(tt.expected)
			if err != nil {
				t.Errorf("could not load expected types file: %v", err)
				return
			}

			// create AST
			schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: string(schemaContent)})
			if err != nil {
				t.Errorf("could not load schema: %v", err)
			}

			query, err := gqlparser.LoadQuery(schema, string(queryContent))
			if err.Error() != "" {
				t.Errorf("could not load query: %v", err)
				return
			}

			// generate types
			types, err := gen.GenerateTypes(context.Background(), schema, query, gen.Options{
				PackageName: "maps",
			})
			if err != nil {
				t.Errorf("could not generate types: %v", err)
				return
			}

			if strings.TrimSpace(string(types)) != strings.TrimSpace(string(expectedTypes)) {
				t.Errorf("expected %s, got %s", string(expectedTypes), string(types))
				return
			}
		})
	}
}
