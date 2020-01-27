package schema

import (
	"html/template"
	"os"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
)

// Generate writes output files for the manifest
func Generate(m *manifest.Manifest) error {

	t, err := template.New("schema").Parse(Template)
	if err != nil {
		return err
	}

	err = t.Execute(os.Stderr, &m)
	if err != nil {
		return err
	}
	return nil
}

// Template is the template definition to create the schema
var Template = `
{{- define "field"}}
	{{ .Name }}
	{{- if not .IsNested -}}
		{{ .ArgList -}}
		:{{" "}}
		{{- if .IsList -}}
			{{.Type }}Connection!
		{{- else -}}
			{{.Type }}
		{{- end -}}
	{{- else -}}
		:{{" "}}
		{{- if .IsList }}[{{ end -}}
		{{ .Type }}{{ if .NonNullable }}!{{ end }}
		{{- if .IsList }}]{{ end -}}
	{{ end }}
{{- end -}}

{{ if .Enums }}
{{range .Enums }}enum {{.Name}} {
	{{range .Values}}{{.}}
	{{end}}
}
{{end}}
{{- end -}}

{{- if .Objects }}
{{range .Objects }}type {{.Name}} {
	{{- range .Fields}}{{ template "field" . }}
	{{- end}}
}
{{end}}

{{- if .Queries}}
type Query {
	{{- range .Queries}}{{ template "field" . }}
	{{- end}}
}
{{end -}}

{{- if .Mutations}}
type Mutation {
	{{- range .Mutations}}{{ template "field" . }}
	{{- end}}
}
{{end -}}
{{ end }}

CONNECTION OBJECTS

FILTER OBJECTS

input TableBooleanFilterInput {
	ne: Boolean
	eq: Boolean
}
input TableIntFilterInput {
	ne: Int
	eq: Int
	le: Int
	lt: Int
	ge: Int
	gt: Int
	contains: Int
	notContains: Int
	between: [Int]
}
input TableStringFilterInput {
	ne: String
	eq: String
	le: String
	lt: String
	ge: String
	gt: String
	contains: String
	notContains: String
	between: [String]
}
input TableFloatFilterInput {
	ne: Float
	eq: Float
	le: Float
	lt: Float
	ge: Float
	gt: Float
	contains: Float
	notContains: Float
	between: [Float]
}
input TableIDFilterInput {
	ne: ID
	eq: ID
	le: ID
	lt: ID
	ge: ID
	gt: ID
	contains: ID
	notContains: ID
	between: [ID]
}
`
