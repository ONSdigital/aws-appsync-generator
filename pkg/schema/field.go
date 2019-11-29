package schema

import (
	"regexp"
)

var reLegalFields = regexp.MustCompile(`(ID|String|Int|Float|Boolean|AWS(Date(Time)?|Time(stamp)?|Email|URL|Phone|IPAddress|JSON))`)

type (
	// Field represents the definition for a graphql object field
	Field struct {
		Name        string
		Type        string
		NonNullable bool `yaml:"nonNullable"`
		Resolver    *Resolver

		// Optional: Type to be used when this field is included in an input
		// object definition if it differs from the main defined type
		InputType string `yaml:"inputType"`
	}

	unmarshalField Field
)

// IsLegalScalarType tests whether a the field is defined
// with one of the allowable graphql (and AWS) scalar types
func (f *Field) IsLegalScalarType() bool {
	return reLegalFields.MatchString(f.Type)
}

// UnmarshalYAML satisfies the custom unmarshaler interface for go-yaml. It's
// called automatically by the unmarshaler to ensure default values get set
// where they haven't been supplied in the user's definition.
func (f *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {

	// Unmarshal to a temporary value otherwise we
	// get broken goroutines for some odd reason
	var u unmarshalField
	if err := unmarshal(&u); err != nil {
		return err
	}

	*f = Field(u)

	if f.Type == "" {
		f.Type = DefaultFieldType
	}
	return nil
}
