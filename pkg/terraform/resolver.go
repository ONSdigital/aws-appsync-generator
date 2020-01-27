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

type datasourceDefinition struct {
	Name string
}

type resolverDefinition struct {
	ParentName string
	FieldName  string
	ArgsSource string

	ParentField *manifest.Field
	Resolver    *manifest.Resolver
	DataSource  *datasourceDefinition
}

func (r resolverDefinition) Write() error {

	name := r.FileName()

	// TODO Add real datasource
	r.DataSource = &datasourceDefinition{
		Name: "dynamo",
	}

	r.ArgsSource = "args"
	// TODO switch to "source" if nested

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

	action := r.Resolver.Action
	if action == "get" && r.ParentField.IsList() {
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

func (r resolverDefinition) KeyFieldMapJSON() string {
	fl := make([]string, len(r.Resolver.KeyFields))
	for i, f := range r.Resolver.KeyFields {
		fl[i] = fmt.Sprintf(`"%s":"%s"`, f.Name, f.Type)
	}
	return "{" + strings.Join(fl, ",") + "}"
}

// Name returns the generated name of the resolver from its
// parent and field
func (r resolverDefinition) Name() string {
	return strings.ToLower(fmt.Sprintf("%s_%s_%s", r.ParentName, r.Resolver.Action, r.FieldName))
}

// FileName returns the generated name of the file the
// terraform will be written to
func (r resolverDefinition) FileName() string {
	// TODO
	// ...
	return strings.ToLower(fmt.Sprintf("%s%s.tf", OutputPrefix, r.Name()))
}

var resolverTemplate = `
## !NOTE: This file is auto-generated DO NOT EDIT
## Generated at {{now}}
resource "aws_appsync_resolver" "{{.Name}}" {
	api_id            = aws_appsync_graphql_api.record.id
	type              = "{{.ParentName}}"
	field             = "{{.FieldName}}"
	data_source       = aws_appsync_datasource.{{ .DataSource.Name }}.name
	request_template  = <<EOF
{{template "request" .}}
EOF
	response_template = <<EOF
{{template "response" .}}
EOF
}
`
