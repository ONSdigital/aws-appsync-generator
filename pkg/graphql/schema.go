package graphql

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// GeneratedFilesPath defines where to put files created from parsing the schema
var GeneratedFilesPath = "./generated"

var schema Schema

type (
	// Schema represents the elements of a graphql schema
	Schema struct {
		Enums   []*Enum   `yaml:"enums"`
		Objects []*Object `yaml:"objects"`
		Queries []*Query  `yaml:"queries"`

		Sources map[string]Source `yaml:"sources"`

		// Automatically populated to create
		// filtering options for list types
		FilterInputs []string

		// Connection objects to be built - populated automatically by "list" resolvers
		Connections   []string
		FilterObjects FilterObjectList

		// Contains any errors raised during the generation process
		Errors []error

		objectLookup map[string]*Object
		// dataSourceType string
	}
)

// NewSchemaFromManifest parses a schema manifest in YAML format and generates
// a new schema struct
func NewSchemaFromManifest(manifest []byte, dataSourceType string) (*Schema, error) {
	var s Schema
	if err := yaml.UnmarshalStrict(manifest, &s); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal schema definition")
	}
	// s.dataSourceType = dataSourceType
	s.FilterInputs = []string{"Int", "String", "Float", "ID"}
	s.Errors = []error{}
	s.Connections = []string{}

	s.objectLookup = make(map[string]*Object)
	for _, o := range s.Objects {
		s.objectLookup[o.Name] = o
	}
	return &s, nil
}

var funcMap = template.FuncMap{
	"now": func() string {
		return time.Now().String()
	},
}

// GenerateBytes renders the schema to a bytes buffer ready to be writen to
// an output stream
func (s *Schema) GenerateBytes() ([]byte, error) {
	generated := bytes.Buffer{}

	if err := schemaTemplate.Execute(&generated, s); err != nil {
		return nil, err
	}
	return generated.Bytes(), nil
}

// CleanOutput empties the output path
func (s *Schema) CleanOutput() error {
	d, err := os.Open(GeneratedFilesPath)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(GeneratedFilesPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func setDataSource(r *Resolver, s *Schema) error {
	if ds, ok := s.Sources[r.SourceKey]; ok {
		r.DataSource = &ds
		return nil
	}

	if ds, ok := s.Sources["default"]; ok {
		r.DataSource = &ds
		return nil
	}

	return fmt.Errorf("resolver '%s_%s' has unknown data source '%s'", r.Parent, r.FieldName, r.SourceKey)
}

// WriteAll outputs the generated public schema and any resolver files to the
// location given by `GeneratedFilesPath`
func (s *Schema) WriteAll() error {

	if err := s.CleanOutput(); err != nil {
		return err
	}

	toWrite := []schemaFileWriter{}

	for _, q := range s.Queries {
		if r := q.Resolver; r != nil {
			r.Parent = "Query"
			r.FieldName = q.Name
			r.ArgsSource = "args"
			if err := setDataSource(r, s); err != nil {
				return err
			}

			// Create appropriate input and connection objects
			if r.Action == ActionList {
				s.Connections = append(s.Connections, r.Type)
				o, ok := s.objectLookup[r.Type]
				if !ok {
					s.Errors = append(s.Errors, fmt.Errorf("unknown type '%s' when attempting to create filter object", r.Type))
					continue
				}
				s.AddFilterFromObject(o)
			}

			toWrite = append(toWrite, r)
		}
	}

	for _, o := range s.Objects {
		for _, f := range o.Fields {
			if r := f.Resolver; r != nil {
				r.Parent = o.Name
				r.FieldName = f.Name
				r.ArgsSource = "source"
				if err := setDataSource(r, s); err != nil {
					return err
				}

				// Create appropriate input and connection objects
				if r.Action == ActionList {
					s.Connections = append(s.Connections, r.Type)
				}

				// TODO - filters

				toWrite = append(toWrite, r)
			}
		}
	}

	// TODO ... mutation resolvers

	// Add the schema to write of course!
	toWrite = append(toWrite, s)

	errored := make(chan error)
	names := make(chan string)

	for _, w := range toWrite {
		go write(w, errored, names)
	}

	for i := 0; i < len(toWrite); i++ {
		select {
		case err := <-errored:
			s.Errors = append(s.Errors, err)
		case name := <-names:
			log.Printf("written: %s", name)
		}
	}

	if len(s.Errors) > 0 {
		return errors.New("errors occured during generation")
	}
	return nil
}

// OutputName returns the file name to be written for the schema
func (s *Schema) OutputName() string {
	return "schema.public.graphql"
}

type schemaFileWriter interface {
	GenerateBytes() ([]byte, error)
	OutputName() string
}

func write(w schemaFileWriter, e chan error, d chan string) {
	bb, err := w.GenerateBytes()
	if err != nil {
		e <- errors.Wrap(err, "failed to generate content")
		return
	}

	path := filepath.Join(GeneratedFilesPath, w.OutputName())

	if err := ioutil.WriteFile(path, bb, 0644); err != nil {
		e <- errors.Wrap(err, "failed to write schema file")
	}
	d <- path
}
