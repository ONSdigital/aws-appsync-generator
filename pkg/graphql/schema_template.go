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
{{.Name}}{{ .Resolver.KeyFieldArgsString }}: {{ if eq .Resolver.Action "get-items" }}[{{end}}{{.Resolver.Type.Name -}}{{ if eq .Resolver.Action "get-items" }}]{{end}}{{ if eq .Resolver.Action "list" }}Connection!{{ end }}
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

{{- range .InputObjects }}
input {{ .Name }} {
	{{ range .Fields -}}{{template "field" .}}
	{{ end -}}
}
{{ end -}}

{{- if .Queries}}
type Query {
	{{- range .Queries}}
	{{ template "resolver" . }}
	{{- end}}
}
{{end -}}

{{- if .Mutations}}
type Mutation {
	{{- range .Mutations}}
	{{ template "resolver" . }}
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
