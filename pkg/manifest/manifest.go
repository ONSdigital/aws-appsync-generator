package manifest

import (
	"gopkg.in/yaml.v2"
)

// Allowable actions for resolvers
const (
	ResolverGet    = "get"
	ResolverInsert = "insert"
	ResolverDelete = "delete"
	ResolverUpdate = "update"
)

type (
	// Manifest is the root of a configuration manifest
	Manifest struct {
		Enums     []Enum   `yaml:"enums"`
		Objects   []Object `yaml:"objects"`
		Queries   []Field  `yaml:"queries"`
		Mutations []Field  `yaml:"mutations"`

		resolvers []*Resolver
	}

	// Enum represents a graphql enumeration definition
	Enum struct {
		Name   string   `yaml:"name"`
		Values []string `yaml:"values"`
	}

	// Object is a graphql object
	Object struct {
		Name   string  `yaml:"name"`
		Fields []Field `yaml:"fields"`
	}

	unmarshalObject Object
)

// UnmarshalYAML custom unmarshals the object field. Sets a flag into its child
// fields to mark them as part of an object field
func (o *Object) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var u unmarshalObject
	if err := unmarshal(&u); err != nil {
		return err
	}

	fields := len(u.Fields)

	o.Name = u.Name
	o.Fields = make([]Field, fields)

	for i := 0; i < fields; i++ {
		f := u.Fields[i]
		f.IsNested = true
		o.Fields[i] = f
	}
	return nil
}

// New creates a new manifest from a yaml configuration
func New(r []byte) (*Manifest, error) {
	var m Manifest
	err := yaml.Unmarshal(r, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
