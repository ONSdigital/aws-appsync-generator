package graphql

import "text/template"

var schemaTemplate = template.Must(template.New("schema").Funcs(funcMap).Parse(`
## !NOTE: This file is auto-generated DO NOT EDIT
## Generated at {{now}}
{{define "field" -}}
{{.Name}}: {{if .Type.IsList}}[{{end -}}
	{{.Type.Name}}{{if .Type.NonNullable}}!{{end}}
{{- if .Type.IsList}}]{{end}}
{{- end}}

{{define "resolver" -}}
{{.Name}}
	{{- if .Resolver.ParamString}}({{.Resolver.ParamString}}){{end -}}
	:
	{{- if eq .Resolver.Action "list"}}[{{end}}
		{{- .Type -}}
	{{- if eq .Resolver.Action "list"}}]{{end}}
{{- end}}

{{- range .Enums}}enum {{.Name}} {
    {{range .Values}}{{.}}
    {{end}}
}
{{end}}

{{- range .Objects}}type {{.Name}} {
    {{range .Fields}}{{template "field" .}}
    {{end}}
}
{{end}}

{{- range .Connections}}
type {{.}}Connection {
	items: [{{.}}]
	nextToken: String
}
{{ end -}}

{{- range .FilterObjects }}
input {{ .Name }} {
	{{ range .Fields -}}{{template "field" .}}
	{{ end -}}
}
{{ end -}}

{{- if .Queries}}
type Query {
	{{- range .Queries}}
	{{.Name}}{{ .Resolver.KeyFieldArgsString }}: {{.Resolver.Type -}}{{ if eq .Resolver.Action "list" }}Connection!{{ end }}
	{{- end}}
}
{{end -}}

input TableBooleanFilterInput {
	ne: Boolean
	eq: Boolean
}
{{range .FilterInputs}}input Table{{.}}FilterInput {
	ne: {{.}}
	eq: {{.}}
	le: {{.}}
	lt: {{.}}
	ge: {{.}}
	gt: {{.}}
	contains: {{.}}
	notContains: {{.}}
	between: [{{.}}]
}
{{end}}
`))
