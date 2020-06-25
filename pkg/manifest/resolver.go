package manifest

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Resolver represents a graphql resolver in the  A resolver can
	// either be a standalone in a Query or Mutation block, or a sub-resolver on
	// an object field
	Resolver struct {
		// Source must be a valid defined `DataSource` name
		DataSourceName string  `yaml:"source"`
		Action         string  `yaml:"action"`
		KeyFields      []Field `yaml:"keyFields"`

		// Custom fields for if the action is specified as "custom"
		Template string `yaml:"template"`

		// Non-parsed calculated fields
		Identifier           string `yaml:"-"`
		ParentType           string `yaml:"-"`
		FieldName            string `yaml:"-"`
		ReturnType           string `yaml:"-"`
		Signature            string `yaml:"-"`
		Request              string `yaml:"-"`
		Response             string `yaml:"-"`
		DataSourceType       string `yaml:"-"`
		DataSourceIdentifier string `yaml:"-"`
	}
)

// UnmarshalYAML performs custom unmarshalling for the Resolver struct. It runs
// simple validation to ensure valid field combinations
func (r *Resolver) UnmarshalYAML(unmarshal func(interface{}) error) error {

	type U Resolver
	var u U

	if err := unmarshal(&u); err != nil {
		return err
	}
	*r = Resolver(u)

	if r.Action == "custom" && r.Template == "" {
		return errors.New("must specify custom template when resolver action is 'custon'")
	}

	if r.Action != "custom" {
		if r.Template != "" {
			return errors.New("must not specify a custom template when resolver action is not 'custom'")
		}
		// For a standard action, the template name is mapped as the
		// name of the action
		r.Template = r.Action
	}

	return nil
}

func (r *Resolver) String() string {
	return fmt.Sprintf("[%s] %s (%s)", r.ParentType, GetAttributeName(r.FieldName), r.Template)
}

// ArgsSource returns "args" or "source" depending on whether the parent is a
// custom object or query/mutation. Expects resolver to have already be populated
// with a valid `ParentType`. If no type is set it will return `source` which
// may not be what you want.
func (r *Resolver) ArgsSource() string {
	switch r.ParentType {
	case "Query":
		return "args"
	case "Mutation":
		return "args"
	default:
		return "source"
	}
}

// KeyFieldJSONMap converts the `KeyFields` into a JSON formatted map suitable
// for use in a VTL template
func (r *Resolver) KeyFieldJSONMap() string {
	fl := make([]string, len(r.KeyFields))
	for i, f := range r.KeyFields {
		// Remove any "mandatory" mark as it doesn't make sense in the usage
		// contexts of the resultant map
		returnType := GetAttributeType(f.Name, "String")
		returnType = strings.TrimSuffix(returnType, "!")
		fl[i] = fmt.Sprintf(`"%s":"%s"`, GetAttributeName(f.Name), returnType)
	}
	return "{" + strings.Join(fl, ",") + "}"
}

// KeyFieldJSONList converts the `KeyFields` into a JSON formatted list of names
// suitable for use in a VTL template
func (r *Resolver) KeyFieldJSONList() string {
	fl := make([]string, len(r.KeyFields))
	for i, f := range r.KeyFields {
		fl[i] = fmt.Sprintf(`"%s"`, GetAttributeName(f.Name))
	}
	return "[" + strings.Join(fl, ",") + "]"
}
