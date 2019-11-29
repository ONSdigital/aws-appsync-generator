package schema

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// DefaultFieldType is the scalar type assigned to fields that don't
// have explicit types set in their definition
const DefaultFieldType = "String"

type (
	// Enum represents a graphql enumeration definition
	Enum struct {
		Name   string
		Values []string
	}

	// Schema represents the whole parsed schema definition
	Schema struct {
		Enums     []*Enum
		Objects   []*Object
		Queries   []*Query
		Inputs    []*Input
		Mutations []*Query

		// Generated
		objectMap map[string]*Object
		enumMap   map[string]*Enum
		inputMap  map[string]*Input
	}
)

// Enum returns the definition for the given object if it exists in the schema
func (s *Schema) Enum(name string) (*Enum, error) {
	if e, ok := s.enumMap[name]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("enum '%s' not found in schema definition", name)
}

// Object returns the definition for the given object if it exists in the schema
func (s *Schema) Object(name string) (*Object, error) {
	if o, ok := s.objectMap[name]; ok {
		return o, nil
	}
	return nil, fmt.Errorf("object '%s' not found in schema definition", name)
}

// Input returns the definition for the given object if it exists in the schema
func (s *Schema) Input(name string) (*Input, error) {
	if i, ok := s.inputMap[name]; ok {
		return i, nil
	}
	return nil, fmt.Errorf("input type '%s' not found in schema definition", name)
}

// MustParse wraps New and ensures a new schema is created.
// It will panic if an error is raised.
func MustParse(definition []byte) *Schema {
	s, err := New(definition)
	if err != nil {
		panic(err)
	}
	return s
}

// New creates a new schema object from a given definition
func New(definition []byte) (*Schema, error) {

	var s Schema
	if err := yaml.UnmarshalStrict(definition, &s); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal schema definition")
	}

	// Generate lookups
	s.enumMap = make(map[string]*Enum)
	for _, e := range s.Enums {
		s.enumMap[e.Name] = e
	}

	s.objectMap = make(map[string]*Object)
	for _, o := range s.Objects {
		o.fieldMap = make(map[string]*Field)
		for _, f := range o.Fields {
			o.fieldMap[f.Name] = f
		}
		s.objectMap[o.Name] = o
	}

	s.inputMap = make(map[string]*Input)
	for _, i := range s.Inputs {
		s.inputMap[i.Name] = i
	}

	return &s, nil
}
