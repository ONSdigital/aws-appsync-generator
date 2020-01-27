package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
)

type resolverData struct {
	// Name of the parent field this resolver is attached to
	FieldName string

	// The direct parent of the field the resolver is attached
	// to. This will either be a user defined object field, or
	// one of Query or Mutation
	ParentObjectName string

	// Action and keyfields are imported from the field's
	// resolver definiton
	Action    string
	KeyFields []manifest.Field

	// The amount of items returned: single, list, paged
	Returns string

	// Defines the type of datasource to connect the resolver to
	DataSourceType string
	DataSourceName string
}

// Name returns the constructed identifer for the resolver
func (r resolverData) Name() string {
	return strings.ToLower(fmt.Sprintf("%s_%s_%s", r.ParentObjectName, r.Action, r.Name()))
}

// IsNested denotes if a resolver is nested if it's a direct
// decendent of a user defined object, and not Query or Mutation.
func (r resolverData) IsNested() bool {
	return r.ParentObjectName != "Query" && r.ParentObjectName != "Mutation"
}

// FileName wraps the constructed resolver name to make it
// appropriate for a filename
func (r resolverData) FileName() string {
	return fmt.Sprintf("%s%s.tf", OutputPrefix, r.Name())
}

func (r resolverData) Write() error {

	name := r.FileName()

	// // TODO Add real datasource
	// r.DataSource = &datasourceDefinition{
	// 	Name: "dynamo",
	// }

	path := filepath.Join(OutputPath, name)

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to generate terraform file")
	}
	defer file.Close()

	t, err := template.New(name).Funcs(funcMap).Parse(resolverTemplate)
	if err != nil {
		return err
	}

	// TODO nested and datasource template parts
	datasource := "dynamo"
	nested := "-nested"
	// TODO

	action := r.Action
	if action == "get" && r.Returns != "single" {
		// If the field return wants a list of values instead
		// a single value, then we update to the virtual action
		// of "list" to get the appropriate template
		action = "list"
	}

	t, err = t.ParseFiles(
		filepath.Join(TemplatePath, "resolvers", datasource, "request", action+nested+".tmpl"),
		filepath.Join(TemplatePath, "resolvers", datasource, "response", action+nested+".tmpl"),
	)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(file, name, r); err != nil {
		return err
	}

	return nil
}
