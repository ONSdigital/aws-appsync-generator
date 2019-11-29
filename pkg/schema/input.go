package schema

import (
	"fmt"

	"github.com/pkg/errors"
)

type (
	// Input defines an input object type
	// These objects are used to define the fields to be passed to a
	// mutation query.
	Input struct {
		// Name gives the name of the input object. By convention it should
		// end with "Input"
		Name string `yaml:"name"`

		// Base defines which standard object to use a template for this input
		// object. Fields from the base are replicated into the input object
		// (with the exception of fields marked as excluded below). The types
		// of the replicated fields are also preserved unless they contain
		// an `inputType` override in their object definition
		Base string `yaml:"base"`

		// Params that should be marked nonNullable in the input defintion.
		// Params listed here must exist in the base object
		NonNullableParams []string `yaml:"nonNullableParams"`

		// Fields from the base object to exclude in the input definition
		// Fields listed here must exist in the base object
		Exclude []string `yaml:"exclude"`

		// Fields define the fields that are included in the import object type.
		// This is usually auto populated from the Base object if specified
		// Fields []*Field `yaml:"fields"`

		Params []*ResolverParam `yaml:"params"` // TODO remove?

		excludeMap   map[string]bool
		nonNullableMap map[string]bool
	}

	unmarshalInput Input
)

func (i *Input) IsExcluded(name string) bool {
	_, ok := i.excludeMap[name]
	return ok
}

func (i *Input) IsNonNullable(name string) bool {
	_, ok := i.nonNullableMap[name]
	return ok
}

// Populate will populate the parameters (fields) for this input object from
// the fields of the base object. Will error if the given base object is not
// found in the provided schema
func (i *Input) Populate(s *Schema) error {
	bo, err := s.Object(i.Base)
	if err != nil {
		return errors.Wrapf(err, "unable to populate '%s', not definition for '%s' in schema", i.Name, i.Base)
	}
	for _, f := range bo.Fields {
		if i.IsExcluded(f.Name) {
			continue
		}

		nonNullable := false
		if i.IsNonNullable(f.Name) {
			nonNullable = true
		}

		// Allow the base object to specify an override
		// value on the parameter type - this allows for 
		// scalar ids to be used to insert where the base
		// object would actually return an object type.
		scalarType := f.Type
		if f.InputType != "" {
			scalarType = f.InputType
		}

		p := &ResolverParam{
			Name:      f.Name,
			Type:      scalarType,
			NonNullable: nonNullable,
		}
		i.Params = append(i.Params, p)
	}
	return nil
}

// UnmarshalYAML satisfies the custom unmarshaler interface for go-yaml. It's
// called automatically by the unmarshaler to ensure default values get set
// where they haven't been supplied in the user's definition.
func (i *Input) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// Unmarshal to a temporary value otherwise we
	// get broken goroutines for some odd reason
	var u unmarshalInput
	if err := unmarshal(&u); err != nil {
		return err
	}

	*i = Input(u)

	// Ensure we've only included EITHER Base object (+nonNullable / exclude lists)
	// or Fields list
	if (i.Base != "" || i.NonNullableParams != nil || i.Exclude != nil) && i.Params != nil {
		return fmt.Errorf("can only specify one of 'base' or 'params' in input '%s'", i.Name)
	}

	if i.Exclude != nil {
		if i.excludeMap == nil {
			i.excludeMap = make(map[string]bool)
		}
		for _, e := range i.Exclude {
			i.excludeMap[e] = true
		}
	}

	if i.NonNullableParams != nil {
		if i.nonNullableMap == nil {
			i.nonNullableMap = make(map[string]bool)
		}
		for _, m := range i.NonNullableParams {
			i.nonNullableMap[m] = true
		}
	}
	return nil
}
