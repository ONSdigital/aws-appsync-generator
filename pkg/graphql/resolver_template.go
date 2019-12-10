package graphql

//template.Must(template.New("resolver").Funcs(funcMap).Parse(`
var resolverTemplate = `
## !NOTE: This file is auto-generated DO NOT EDIT
## Generated at {{now}}
resource "aws_appsync_resolver" "{{.Parent}}_{{.FieldName}}" {
	api_id            = aws_appsync_graphql_api.record.id
	type              = "{{.Parent}}"
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

//))
