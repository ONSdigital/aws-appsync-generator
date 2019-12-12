package graphql

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Constants for resolver actions
const (
	ActionGet    = "get"
	ActionList   = "list"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionInsert = "insert"
)

type (
	// Resolver is the represenatation of a field resolver
	Resolver struct {
		// The Action for the resolver to perform - get, list, insert, update, delete
		Action string `yaml:"action"`

		// The graphql type that the resolver returns. Where this resolver action
		// is other than get or delete, set this to be the base object type.
		// Appropriate "Input" variants of the object will then be defined.
		// Where the action is "list", the base type will also have an associated
		// "Connection" object defined.
		Type string `yaml:"type"`

		// Set of fields to be used as keys in the query (primary key etc)
		KeyFields []*Field `yaml:"keyFields"`

		// Name of the datasource to be used for this resolver. Must have
		// been defined in the `sources` section of the manifest
		SourceKey string `yaml:"source"`

		// The below are set automatically as the schema is parsed. They should
		// not be included in the manifest YAML.
		DataSource *Source // Key to a datasource defined in the manifest
		ArgsSource string  // $ctx.{{ArgSource}}.get() - "args" or "source"
		Parent     string  // Parent field
		FieldName  string  // Field name attached to
	}
)

// KeyFieldJSONMap converts the `KeyFields` into a JSON formatted map
func (r *Resolver) KeyFieldJSONMap() string {
	fl := make([]string, len(r.KeyFields))
	for i, f := range r.KeyFields {
		fl[i] = fmt.Sprintf(`"%s":"%s"`, f.Name, f.Type.Name)
	}
	return "{" + strings.Join(fl, ",") + "}"
}

// KeyFieldJSONList converts the `KeyFields` into a JSON formatted list of names
func (r *Resolver) KeyFieldJSONList() string {
	fl := make([]string, len(r.KeyFields))
	for i, f := range r.KeyFields {
		fl[i] = fmt.Sprintf(`"%s"`, f.Name)
	}
	return "[" + strings.Join(fl, ",") + "]"
}

// KeyFieldArgsString returns the keyfield names and types in a string format
// suitable to be used as the arguments list in a resolver defintion in
// the schemea
func (r *Resolver) KeyFieldArgsString() string {

	switch r.Action {
	case ActionList:
		return fmt.Sprintf("(filter: %sFilter, limit: Int, nextToken: String)", r.Type)
	case ActionInsert:
		return fmt.Sprintf("(input: Create%sInput)", r.Type)
	case ActionUpdate:
		return fmt.Sprintf("(input: Update%sInput)", r.Type)
	}

	if l := len(r.KeyFields); l > 0 {
		fl := make([]string, l)
		for i, f := range r.KeyFields {
			fl[i] = fmt.Sprintf(`%s:%s`, f.Name, f.Type.Name)
		}
		return "(" + strings.Join(fl, ",") + ")"
	}
	return ""
}

// GenerateBytes renders the resolver ready to be written to an output stream
func (r *Resolver) GenerateBytes() ([]byte, error) {
	generated := bytes.Buffer{}

	t, err := template.New(r.Action).Funcs(funcMap).Parse(resolverTemplate)
	if err != nil {
		return nil, err
	}

	t, err = t.ParseFiles(
		"templates/resolvers/dynamo/request/"+r.Action+".tmpl",
		"templates/resolvers/dynamo/response/"+r.Action+".tmpl",
	)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(&generated, r); err != nil {
		return nil, err
	}
	return generated.Bytes(), nil
}

// OutputName returns the file name to be written for the resolver
func (r *Resolver) OutputName() string {
	return strings.ToLower(fmt.Sprintf("_%s_%s.tf", r.Parent, r.FieldName))
}
