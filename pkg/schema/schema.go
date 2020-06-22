package schema

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
)

type (

	// Schema represents a graphql schema
	Schema struct {
		Objects map[string]Object
		Enums   map[string]Enum
		Queries map[string]Query
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

	Query struct {
		Name string
	}
)

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

}
{{- end }}
type Schema {
	Query: Query
}

####### DELETE THIS!!!
type Query {
	food(id:ID!): Food
}
`))

// New creates an empty schema
func New() *Schema {
	return &Schema{
		Objects: make(map[string]Object),
		Enums:   make(map[string]Enum),
	}
}

// NewFromManifest creates a new schema from a parsed manifest
func NewFromManifest(m *manifest.Manifest) (*Schema, error) {

	s := New()
	s.importEnumsFromManifest(m)
	s.importObjectsFromManifest(m)

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
