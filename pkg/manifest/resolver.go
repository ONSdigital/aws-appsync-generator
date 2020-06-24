package manifest

import (
	"github.com/pkg/errors"
)

type (
	// Resolver represents a graphql resolver in the manifest. A resolver can
	// either be a standalone in a Query or Mutation block, or a sub-resolver on
	// an object field
	Resolver struct {
		// Source must be a valid defined `DataSource` name
		Source    string  `yaml:"source"`
		Action    string  `yaml:"action"`
		KeyFields []Field `yaml:"keyFields"`

		// Custom fields for if the action is specified as "custom"
		Template string `yaml:"template"`

		// Non-parsed fields
		ParentType string `yaml:"-"`
	}
)

// UnmarshalYAML performs custom unmarshalling for the Resolver struct. It runs
// simple validation to ensure valid field combinations
func (m *Resolver) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type U Resolver
	var u U

	if err := unmarshal(&u); err != nil {
		return err
	}
	*m = Resolver(u)

	if m.Action == "custom" && m.Template == "" {
		return errors.New("must specify custom template when resolver action is 'custon'")
	}

	if m.Action != "custom" {
		if m.Template != "" {
			return errors.New("must not specify a custom template when resolver action is not 'custom'")
		}
		// For a standard action, the template name is mapped as the
		// name of the action
		m.Template = m.Action
	}

	return nil
}
