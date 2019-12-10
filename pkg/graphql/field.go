package graphql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var reLegalFields = regexp.MustCompile(`(ID|String|Int|Float|Boolean|AWS(Date(Time)?|Time(stamp)?|Email|URL|Phone|IPAddress|JSON))`)

const defaultFieldType = "String"

// Custom errors
var (
	ErrFieldHasNoName            = errors.New("fields must have a 'name' attribute")
	ErrFieldHasBadTypeDefinition = errors.New("field has bad type definition, must be <type> or [<type>]")
	ErrTypeAndResolver           = errors.New("field cannot declare Type and Resolver")
)

type (
	// FieldType represents the type associated with a field.
	FieldType struct {
		Name        string
		IsList      bool
		NonNullable bool
	}

	// Field represents a field of a graphql object
	Field struct {
		Name string

		// (Optional) Type defines the scalar or object type of the field. If a type is not
		// specified, it will default to the `String` scalar type if not set and
		// Resolver has not been specified
		// Must specify exactly one of Type OR Resolver
		Type *FieldType

		// (Optional) Define a resolver to fetch the value for this field
		// Resolver *Resolver `yaml:"resolver,omitempty"`
		// Must specify exactly one of Type OR Resolver
		Resolver *Resolver

		// (Optional) Type to be used when this field is included in an input
		// object definition if it differs from the main defined type
		InputType *FieldType
	}
)

// IsLegalScalarType tests whether a the field is defined
// with one of the allowable graphql (and AWS) scalar types
func (f *Field) IsLegalScalarType() bool {
	return reLegalFields.MatchString(f.Type.Name)
}

// UnmarshalYAML satisfies the custom unmarshaler interface for go-yaml. It is
// called automatically by the YAML unmarshal to deal with the [] list notation in
// field type definitions as they need to be read as scalar values, not YAML lists.
func (ft *FieldType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var u interface{}
	if err := unmarshal(&u); err != nil {
		return err
	}

	// Some ugly unboxing. We need to do this as the manifest
	// may try to specify a [type] for a field which needs to
	// become a list of that type but needs to co-exist with
	// other fields where it's not a list... Yaml :)
	switch t := u.(type) {
	case nil:
		ft.Name = defaultFieldType
	case string:
		ft.Name = fmt.Sprintf("%s", u)
	case []interface{}:
		switch tl := t[0].(type) {
		case string:
			ft.Name = fmt.Sprintf("%s", tl)
			ft.IsList = true
		}
	default:
		return ErrFieldHasBadTypeDefinition
	}
	ft.parseNonNullability()

	return nil
}

func (ft *FieldType) parseNonNullability() {
	if strings.HasSuffix(ft.Name, "!") {
		ft.Name = ft.Name[:len(ft.Name)-1]
		ft.NonNullable = true
	}
}

// UnmarshalYAML satisfies the custom unmarshaler interface for go-yaml. It's
// called automatically by the unmarshaler to ensure default values get set
// where they haven't been supplied in the user's definition.
func (f *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// Use a custom struct to temporarily marshal into so we
	// have control over type checking and defauting fields.
	var u struct {
		Name      string     `yaml:"name"`
		Resolver  *Resolver  `yaml:"resolver"`
		Type      *FieldType `yaml:"type"`
		InputType *FieldType `yaml:"inputType"`
	}
	if err := unmarshal(&u); err != nil {
		return err
	}

	// Map the fields from the temp to the real struct
	// NOTE! If the fields allowable change then remember
	// to update this mapping!
	f.Name = u.Name
	f.Resolver = u.Resolver
	f.Type = u.Type
	f.InputType = u.InputType

	// Ensure at least ONE of Type or Resolver is set (and default
	// to a Type:String if neither)
	if u.Resolver != nil && u.Type != nil {
		return ErrTypeAndResolver
	}

	// If the type isn't set and there is no resolver,
	// then type needs to get the default type
	if u.Resolver == nil && u.Type == nil {
		f.Type = &FieldType{
			Name:   defaultFieldType,
			IsList: false,
		}
	}

	return nil
}
