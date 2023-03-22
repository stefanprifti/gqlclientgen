package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/stefanprifti/gql/gen"
	"github.com/stefanprifti/gql/introspect"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

//go:embed client.go.tmpl
var clientFileTmpl string

const (
	clientFile        = "client.go"
	modelFile         = "model.go"
	gqlSchemaFile     = "schema.graphql"
	gqlIntrospectFile = "schema.introspect.json"
	configFile        = "gqlclientgen.yml"
)

type ClientMethod struct {
	Name     string
	Query    string
	Request  string
	Response string
	Type     string
}

type Operation struct {
	FilePath    string
	FileContent []byte
	GQLTypes    []byte
	Doc         *ast.QueryDocument
}

type Service struct {
	Package string

	SchemaURL string
	// SchemaContent is the schema in GQL format
	SchemaContent string
	// SchemaDoc is the schema in AST format
	SchemaDoc *ast.Schema
	// SchemaJSON is the schema in JSON format
	SchemaJSON []byte

	OperationsFolder string
	OperationDocs    []Operation

	ClientFolder string
}

type App struct {
	Config   Config
	Services []Service
}

func New(config Config) (*App, error) {
	app := &App{}

	err := app.setConfig(config)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) setConfig(config Config) error {
	services := make([]Service, 0, len(config.Services))

	for _, service := range config.Services {
		services = append(services, Service{
			Package:          service.Package,
			SchemaURL:        service.URL,
			OperationsFolder: service.Operations.Root,
			ClientFolder:     service.Client.Root,
		})
	}

	a.Services = services

	return nil
}

func (s *Service) ResolveSchema() error {
	schema, err := FetchSchema(s.SchemaURL)
	if err != nil {
		return fmt.Errorf("failed to fetch schema: %w", err)
	}

	schemaBytes, err := introspect.SchemaToText(schema)
	if err != nil {
		return fmt.Errorf("failed to convert schema to text: %w", err)
	}

	doc, err := gqlparser.LoadSchema(&ast.Source{Name: fmt.Sprintf("%s.schema.graphql", s.Package), Input: string(schemaBytes)})
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	s.SchemaContent = string(schemaBytes)
	s.SchemaDoc = doc
	s.SchemaJSON, err = json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	return nil
}

func (s *Service) ResolveOperations() error {
	// read filees in a folder
	files, err := os.ReadDir(s.OperationsFolder)
	if err != nil {
		return fmt.Errorf("failed to read folder %s: %w", s.OperationsFolder, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".graphql" {
			continue
		}

		// read the file
		f, err := os.OpenFile(filepath.Join(s.OperationsFolder, file.Name()), os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file.Name(), err)
		}

		body, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
		}

		operationDoc, err := gqlparser.LoadQuery(s.SchemaDoc, string(body))
		if err.Error() != "" {
			return fmt.Errorf("failed to load query %s: %w", file.Name(), err)
		}

		s.OperationDocs = append(s.OperationDocs, Operation{
			FilePath:    filepath.Join(s.OperationsFolder, file.Name()),
			FileContent: body,
			Doc:         operationDoc,
		})
	}

	return nil
}

// GenerateIntrospectionFile generates the introspection file for the service
func (s *Service) GenerateIntrospectionFile() error {
	introspectFilePath := filepath.Join(s.ClientFolder, gqlIntrospectFile)
	err := writeFile(introspectFilePath, s.SchemaJSON)
	if err != nil {
		return fmt.Errorf("failed to write introspection file: %w", err)
	}
	return nil
}

// GenerateSchemaFile generates the schema file for the service
func (s *Service) GenerateSchemaFile() error {
	schemaFilePath := filepath.Join(s.ClientFolder, gqlSchemaFile)
	err := writeFile(schemaFilePath, []byte(s.SchemaContent))
	if err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	return nil
}

// GenerateModelFile generates the model file for the service
func (s *Service) GenerateModelFile() error {
	var b bytes.Buffer

	schemaTypes, err := gen.GenerateTypesFromSchema(s.Package, s.SchemaDoc)
	if err != nil {
		return fmt.Errorf("failed to generate types from schema: %w", err)
	}

	b.Write(schemaTypes)

	for _, operation := range s.OperationDocs {
		operationTypes := gen.GenerateTypesFromOperation(operation.Doc)
		b.Write(operationTypes)
	}

	modelFilePath := filepath.Join(s.ClientFolder, modelFile)
	err = writeFile(modelFilePath, b.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	return nil
}

// GenerateClientFile generates the client file for the service
func (s *Service) GenerateClientFile() error {
	// Get the template
	tmpl := template.Must(template.New("template").Parse(clientFileTmpl))

	// Create data for template
	methods := make([]ClientMethod, 0, len(s.OperationDocs))

	for _, operation := range s.OperationDocs {
		methods = append(methods, ClientMethod{
			Name:     operation.Doc.Operations[0].Name,
			Query:    string(operation.FileContent),
			Request:  fmt.Sprintf("%sRequest", operation.Doc.Operations[0].Name),
			Response: fmt.Sprintf("%sResponse", operation.Doc.Operations[0].Name),
			Type:     string(operation.Doc.Operations[0].Operation),
		})
	}

	// Execute the template
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, map[string]interface{}{
		"PackageName": s.Package,
		"Methods":     methods,
	})
	if err != nil {
		return err
	}

	// Get the generated code
	generatedClientCode := buf.String()

	clientFilePath := filepath.Join(s.ClientFolder, clientFile)
	err = writeFile(clientFilePath, []byte(generatedClientCode))
	if err != nil {
		return fmt.Errorf("failed to write client file: %w", err)
	}

	return nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from panic: ", r)
		}
	}()

	config, err := LoadConfig(configFile)
	if err != nil {
		if errors.Is(err, ErrConfigNotFound) {
			fmt.Println("config file gqlclientgen.yml not found")
			return
		}
	}

	app, err := New(config)
	if err != nil {
		fmt.Println("could not initialize app: ", err)
		return
	}

	err = run(app)
	if err != nil {
		fmt.Println("could not generate client: ", err)
		return
	}

	fmt.Println("done")
}

func run(app *App) error {
	for _, service := range app.Services {
		fmt.Println("processing service: ", service.Package, " at ", service.SchemaURL, " with operations at ", service.OperationsFolder, " and client at ", service.ClientFolder, "")

		err := service.ResolveSchema()
		if err != nil {
			return err
		}

		err = service.ResolveOperations()
		if err != nil {
			return err
		}

		err = service.GenerateIntrospectionFile()
		if err != nil {
			return err
		}

		err = service.GenerateSchemaFile()
		if err != nil {
			return err
		}

		err = service.GenerateModelFile()
		if err != nil {
			return err
		}

		err = service.GenerateClientFile()
		if err != nil {
			return err
		}

		fmt.Println("Successfully generated client for service: ", service.Package)
	}

	return nil
}

func writeFile(filePath string, body []byte) error {
	f, err := openFile(filePath)
	if err != nil {
		return err
	}

	// format file content
	if strings.HasSuffix(filePath, ".go") {
		body, err = format.Source(body)
		if err != nil {
			return err
		}
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	return nil
}

// WriteSchema writes the schema to the given writer
func openFile(filePath string) (*os.File, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// create the file even if the directories don't exist
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return nil, err
			}

			f, err = os.Create(filePath)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return f, nil
}

// FetchSchema fetches the schema from the given url
func FetchSchema(url string) (*introspect.Schema, error) {
	schema, err := introspect.URL(url)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
