package schema

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// type Resolver interface{
// 	AppsyncFuzzyMapString() string
// 	AppsyncFieldsMapString() string
// }

type (
	// ResolverParam is a parameter required by a resolver to execute a query
	ResolverParam struct {
		// Name of the field
		Name string

		// Used when a sub field of an object to define
		// which parent field should be used to match to.
		SourceName string `yaml:"sourceName"`

		// The scalar or object type of the parameter
		Type string

		// Denotes whether to do a LIKE match or an exact match in queries
		Fuzzy bool

		// Denotes whether this parameter should be
		// in the resolver signature mandatory
		NonNullable bool `yaml:"nonNullable"`
	}

	unmarshalResolverParam ResolverParam
)

// UnmarshalYAML is called automatically by `yaml.Unmarshal` and is used to
// set defaults into the unmarshalled struct
func (r *ResolverParam) UnmarshalYAML(unmarshal func(interface{}) error) error {
	u := unmarshalResolverParam{}
	if err := unmarshal(&u); err != nil {
		return errors.Wrap(err, "bad resolver config")
	}

	*r = ResolverParam(u)
	if r.Type == "" {
		r.Type = DefaultFieldType
	}
	return nil
}

type (
	// Resolver represents the configuration for a field resolver
	Resolver struct {
		Params []*ResolverParam

		// Gives the execution type of the resolver - get, list, insert
		// or update. Unless overriden also denotes the templates to be
		// used to generate terraform ouput
		// Where the type is 'list", the output type is assumed to be
		// a [list]
		Action string

		// Allow the definition of custom templates if we wish to differ from
		// the standard default set. If not supplied, then both templates used
		// will default the the "Type" above
		CustomRequestTemplateName  string `yaml:"requestTemplate"`
		CustomResponseTemplateName string `yaml:"responseTemplate"`

		// KeyFields (only applicable to resolvers for mutation queries) are
		// the identifiers to be used to pull back the mutated record
		KeyFields []*Field `yaml:"keyFields"`

		// Extra arbitary data to be interpolated into the templates
		Data map[string]interface{}

		// Used for mutations to define the input object
		// to be used for data entry
		Input string

		keyFieldsString string
	}

	// UnmarshalResolver is used as an intermediary for
	// unmarshaling the YAML definition of the field
	unmarshalResolver Resolver
)

func (r *Resolver) hasLegalResolverAction() bool {
	return regexp.MustCompile(`get|list|update|insert|delete`).MatchString(r.Action)
}

// TemplateNames returns the names of the templates
// to be used to generare this resolver
func (r *Resolver) TemplateNames() (string, string, error) {
	if !r.hasLegalResolverAction() {
		return "", "", fmt.Errorf("action type '%s' not allowable for resolver", r.Action)
	}
	if r.Action == "custom" {
		if r.CustomRequestTemplateName == "" || r.CustomResponseTemplateName == "" {
			return "", "", errors.New("action for resolver defined as 'custom' but one or more of custom request/response templates have not been defined")
		}
		return r.CustomRequestTemplateName, r.CustomResponseTemplateName, nil
	}
	return r.Action + ".tmpl", r.Action + ".tmpl", nil
}

// ParamString returns the generated parameter string for a resolver
func (r *Resolver) ParamString() string {
	return r.joinParams(", ", false)
}

// FieldsString returns the fields names list that we pass to the template. This is a
// comma separated list of field names
func (r *Resolver) FieldsString() string {
	return r.joinParams(",", true)
}

// AppsyncFuzzyMapString returns a JSON map containing a mapping of the params in
// the mutation to whether the should match on direct (=) or fuzzy (LIKE) operators
func (r *Resolver) AppsyncFuzzyMapString() string {
	s := make([]string, 0, len(r.Params))
	for _, p := range r.Params {
		s = append(s, fmt.Sprintf(`"%s": %v`, p.Name, p.Fuzzy))
	}
	return "{" + strings.Join(s, ",") + "}"
}

// AppsyncFieldsMapString returns a string representation of an appsync
// compatible map type
func (r *Resolver) AppsyncFieldsMapString() string {
	s := make([]string, 0, len(r.Params))
	for _, p := range r.Params {
		source := p.Name
		if p.SourceName != "" {
			source = p.SourceName
		}
		s = append(s, fmt.Sprintf(`"%s": "%s"`, p.Name, source))
	}
	return "{" + strings.Join(s, ",") + "}"
}

func (r *Resolver) joinParams(joiner string, namesOnly bool) string {
	list := make([]string, len(r.Params))
	for i, p := range r.Params {

		if namesOnly {
			list[i] = p.Name
			continue
		}

		t := p.Type
		if p.NonNullable {
			t = t + "!"
		}
		list[i] = fmt.Sprintf("%s: %s", p.Name, t)
	}
	return strings.Join(list, joiner)
}

func (r *Resolver) KeyFieldsString() string {
	names := make([]string, 0, len(r.KeyFields))
	for _, k := range r.KeyFields {
		names = append(names, `"`+k.Name+`"`)
	}
	return "[" + strings.Join(names, ",") + "]"
}
