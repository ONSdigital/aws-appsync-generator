package manifest

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type (
	// Field is a graphql object field
	Field struct {
		// Marshaled fields
		Name     string    `yaml:"name"`
		Type     string    `yaml:"type"`
		Resolver *Resolver `yaml:"resolver,omitempty"`
		// DataSource *DataSource `yaml:"source"` // TODO define data source

		// Generated fields
		ArgList     string `yaml:"-"`
		IsNested    bool   `yaml:"-"`
		IsList      bool   `yaml:"_"`
		NonNullable bool   `yaml:"_"`

		parentName    string
		isObjectField bool
		// isList        bool
	}

	// Resolver is a grapql field resolver
	Resolver struct {
		Action    string  `yaml:"action"`
		KeyFields []Field `yaml:"keyFields"`
		Paged     bool    `yaml:"paged"`

		// Generated
		ArgList string

		returns string
	}
)

// Field errors
var (
	ErrMissingMandatoryValues    = errors.New("field is missing one or more mandatory fields")
	ErrFieldHasBadTypeDefinition = errors.New("field has bad type definition, must be <type> or [<type>]")
	ErrResolverMustSpecifyAction = errors.New("resolver must specify a valid action type")
	ErrBadResolverAction         = errors.New("resolver specifies invalid action")
)

var reValidResolverActions = regexp.MustCompile(ResolverGet + "|" + ResolverDelete + "|" + ResolverInsert + "|" + ResolverUpdate)

// Returns gives the type of return value for this resolver: single, list or paged
func (r *Resolver) Returns() string {
	return r.returns
}

// UnmarshalYAML custom unmarshals the field and does mandatory value checking
func (f *Field) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var u struct {
		Name     string      `yaml:"name"`
		Type     interface{} `yaml:"type"` // Will be type-asserted for list/non nullability
		Resolver *Resolver   `yaml:"resolver"`
	}

	if err := unmarshal(&u); err != nil {
		return err
	}

	// Check mandatory fields that must be set in the resolver definition
	if u.Name == "" {
		return errors.Wrap(ErrMissingMandatoryValues, "missing 'name'")
	}
	f.Name = u.Name

	// Parse the type to determine if it's a standard scalar type, a user
	// defined object type, or a list varient of those. If it's not set, then
	// it will be defaulted to the graphql.String type
	switch t := u.Type.(type) {
	case nil:
		f.Type = "String"
	case string:
		f.Type = t
	case []interface{}:
		switch tl := t[0].(type) {
		case string:
			f.Type = tl
			f.IsList = true
		}
	default:
		return ErrFieldHasBadTypeDefinition
	}

	// Parse non-nullability from the field
	if strings.HasSuffix(f.Type, "!") {
		f.Type = f.Type[:len(f.Type)-1]
		f.NonNullable = true
	}

	if u.Resolver != nil {
		r := u.Resolver

		if r.Action == "" {
			return errors.Wrapf(ErrResolverMustSpecifyAction, "invalid field definition '%s'", u.Name)
		}

		if !reValidResolverActions.MatchString(r.Action) {
			return errors.Wrapf(ErrBadResolverAction, "invalid action: "+r.Action)
		}

		// Generate the arg list for the field signature if a Query or Mutation
		if !f.IsNested {
			switch r.Action {
			case "get":
				if f.IsList {
					f.ArgList = fmt.Sprintf("(filter: %sFilter, limit: Int, nextToken: String)", f.Type)
					break
				}
				if l := len(r.KeyFields); l > 0 {
					fl := make([]string, l)
					for i, f := range r.KeyFields {
						fl[i] = fmt.Sprintf(`%s:%s`, f.Name, f.Type)
						if f.NonNullable {
							fl[i] += "!"
						}
					}
					f.ArgList = "(" + strings.Join(fl, ",") + ")"
				}
			case "insert":
				f.ArgList = fmt.Sprintf("(input: Create%sInput)", f.Type)
			case "update":
				f.ArgList = fmt.Sprintf("(input: Update%sInput)", f.Type)
			default:
				return errors.New("Unsupported action " + r.Action)
			}
		}

		// If the resolver should returned paged results, then the Type is updated
		// to be a Connection type
		// This will also cause the associated Connection object to be generated
		// TODO Generate connection object

		f.Resolver = r
	}

	return nil
}

// KeyFieldArgsString returns the keyfield names and types in a string format
// suitable to be used as the arguments list in a resolver defintion in
// the schemea
// func (r *Resolver) KeyFieldArgsString() string {

// 	switch r.Action {
// 	case ActionList:
// 		return fmt.Sprintf("(filter: %sFilter, limit: Int, nextToken: String)", r.Type.Name)
// 	case ActionInsert:
// 		return fmt.Sprintf("(input: Create%sInput)", r.Type.Name)
// 	case ActionUpdate:
// 		return fmt.Sprintf("(input: Update%sInput)", r.Type.Name)
// 	}

// 	if l := len(r.KeyFields); l > 0 {
// 		fl := make([]string, l)
// 		for i, f := range r.KeyFields {
// 			fl[i] = fmt.Sprintf(`%s:%s`, f.Name, f.Type.Name)
// 		}
// 		return "(" + strings.Join(fl, ",") + ")"
// 	}
// 	return ""
// }

// IsList returns if cardinality of the field is a list
// func (f *Field) IsList() bool {
// 	return f.isList
// }

// NonNullable returns whether this field is marked as non nullable
// func (f *Field) NonNullable() bool {
// 	return f.nonNullable
// }

// SetIsObjectField marks whether this field is a direct descendent of
// a user defined graphql object type
// func (f *Field) SetIsObjectField(is bool) {
// 	f.isObjectField = is
// }

// ResolverArgList returns schema formatted string suitable to be inserted as an
// arg list for a resolver
func (f *Field) ResolverArgList() string {
	return "" // TODO
}

// IsNotObjectField returns true is this field is NOT a direct child of a custom
// object type. Used when generating the schema to knoew whether
// func (f *Field) IsNotObjectField() bool {
// 	return !f.isObjectField
// }
