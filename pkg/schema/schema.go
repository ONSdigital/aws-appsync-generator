package schema

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"text/template"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/ONSdigital/aws-appsync-generator/pkg/mapping"
	"github.com/pkg/errors"
)

type (

	// Schema represents a graphql schema
	Schema struct {
		Objects   map[string]Object
		Enums     map[string]Enum
		Queries   map[string]*Query
		Mutations map[string]*Query

		resolversByParent manifest.ResolverMap
	}

	// Object represents a graphql object type
	Object struct {
		Name   string
		Fields []string
	}

	// Enum represents a graphql enum type
	Enum struct {
		Name   string
		Values []string
	}

	// Query represents a query (or a mutation - they're the same as far
	// as the schema representation is concerned)
	Query struct {
		Name       string
		ReturnType string
		Signature  string
	}
)

// String stringifies a query suitable for output into a graphql schema
func (q *Query) String() string {
	if q.Signature == "" {
		return fmt.Sprintf("%s: %s", q.Name, q.ReturnType)
	}
	return fmt.Sprintf("%s(%s): %s", q.Name, q.Signature, q.ReturnType)
}

var tmplSchema = template.Must(template.New("schema").Parse(`
{{ range .Enums }}
enum {{ .Name }} {
{{ range .Values }}    {{ . }}
{{ end -}}
}
{{ end }}

{{- range .Objects }}
type {{ .Name }} {
{{ range .Fields }}    {{ . }}
{{ end -}}
}
{{ end }}

{{- if .Queries -}}
type Query {
{{ range .Queries }}    {{ . }}
{{ end -}}
}
{{- end }}
type Schema {
	Query: Query
}

`))

// New creates an empty schema
func New() *Schema {
	return &Schema{
		Objects:   make(map[string]Object),
		Enums:     make(map[string]Enum),
		Queries:   make(map[string]*Query),
		Mutations: make(map[string]*Query),
	}
}

// NewFromManifest creates a new schema from a parsed manifest
func NewFromManifest(m *manifest.Manifest, templates mapping.Templates) (*Schema, error) {

	var err error

	s := New()

	s.resolversByParent, err = m.Resolvers()
	if err != nil {
		return nil, err
	}

	s.importEnumsFromManifest(m)
	s.importObjectsFromManifest(m)

	if err := s.importQueriesFromManifest(m, templates); err != nil {
		return nil, err
	}

	// TODO Import mutations

	// TODO Write associated objects (filters and inputs as required)

	return s, nil
}

func (s *Schema) importEnumsFromManifest(m *manifest.Manifest) {
	for name, values := range m.Enums {
		e := Enum{
			Name:   name,
			Values: make([]string, len(values)),
		}
		copy(e.Values, values)
		s.Enums[name] = e
	}
	return
}

func (s *Schema) importQueriesFromManifest(m *manifest.Manifest, templates mapping.Templates) error {
	for _, f := range m.Queries {
		// Get the datasource for the associated resolver
		// and determine the type
		dataSourceName := f.Resolver.DataSourceName
		dataSource := m.DataSources[dataSourceName]

		dataSourceType := ""
		switch {
		case dataSource.Dynamo != nil:
			dataSourceType = "dynamo"
		default:
			log.Fatalf("unsupported data type in schema: %s", dataSourceName)
		}

		// Get the mapping for resolver and datasource type
		template, err := templates.Get(dataSourceType, f.Resolver.Template)
		if err != nil {
			return err
		}

		fieldName := manifest.GetAttributeName(f.Name)
		returnType := manifest.GetAttributeTypeStripped(f.Name, "")

		r := f.Resolver
		r.Identifier = strings.ToLower(fmt.Sprintf("%s_%s_%s", dataSourceType, r.ParentType, fieldName))
		r.ReturnType = returnType
		r.FieldName = fieldName
		r.DataSourceType = dataSourceType

		// Need to execute the template to get the signature
		var b bytes.Buffer
		if err := template.Template.ExecuteTemplate(&b, "signature", r); err != nil {
			return err
		}
		signature := b.String()
		b.Reset()

		queryName := manifest.GetAttributeName(f.Name)

		s.Queries[queryName] = &Query{
			Name:       queryName,
			ReturnType: manifest.GetAttributeType(f.Name, ""),
			Signature:  signature,
		}

	}
	return nil
}

func (s *Schema) importObjectsFromManifest(m *manifest.Manifest) {
	for name, fields := range m.Objects {
		numFields := len(fields)

		// Create the baseline object
		o := Object{
			Name:   name,
			Fields: make([]string, numFields),
		}

		// Process each field for correct output format. Default each type to
		// a graphql String if no other type has been specified.
		for f := 0; f < numFields; f++ {
			fieldDefinition := fields[f].Name
			if !strings.Contains(fieldDefinition, ":") {
				fieldDefinition += ":String"
			}
			o.Fields[f] = strings.ReplaceAll(fieldDefinition, ":", ": ")
		}

		// Store the parsed object into the schema
		s.Objects[name] = o
	}
	return
}

// Marshal marshals the schema into bytes
func (s *Schema) Marshal() (out []byte, err error) {
	bb := &bytes.Buffer{}
	if err := tmplSchema.Execute(bb, s); err != nil {
		return nil, err
	}
	out = bb.Bytes()
	return
}

// Write outouts the generated terraform to a given io.Writer
func (s *Schema) Write(w io.Writer) error {
	data, err := s.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to write schema")
	}
	_, err = w.Write(data)
	return err
}
