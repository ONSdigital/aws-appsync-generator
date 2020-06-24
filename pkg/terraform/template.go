package terraform

import "text/template"

var tmplTerraform = template.Must(template.New("terraform").Parse(`
# ===============================================
# APPSYNC API
# ===============================================
{{- with .AppSync}}
resource "aws_iam_role" "appsync" {
	name = "${terraform.workspace}-{{ .Name }}"
	assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Service\":\"appsync.amazonaws.com\"},\"Effect\":\"Allow\"}]}"
}
resource "aws_iam_role_policy_attachment" "appsync" {
	policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs"
	role       = aws_iam_role.appsync.name
}
resource "aws_appsync_graphql_api" "appsync" {
	authentication_type = "API_KEY"
	name                = "${terraform.workspace}-{{ .Name }}"
	schema              = "${file("schema.public.graphql")}"
	log_config {
		cloudwatch_logs_role_arn = aws_iam_role.appsync.arn
		field_log_level          = "ERROR"
	}
	tags = {
		Environment = terraform.workspace
		Deployment = "${terraform.workspace}-{{ .Name }}"
	}
}
output "graphql_api_id" {
	value = aws_appsync_graphql_api.appsync.id
}
output "graphql_host" {
	value = aws_appsync_graphql_api.appsync.uris
}
{{end}}
# ===============================================
# DATA SOURCES
# ===============================================
{{range .DataSources.Dynamo}}
# Dynamo: {{ .Name }} ----
resource "aws_iam_role_policy" "record_dynamo_{{ .Identifier }}" {
	name		= "${terraform.workspace}-dynamo-{{.Name}}"
	role 		= aws_iam_role.appsync.id
	policy 		= "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"dynamodb:*\"],\"Effect\":\"Allow\",\"Resource\":[\"${aws_dynamodb_table.{{ .Identifier }}.arn}\"]}]}"
}
resource "aws_dynamodb_table" "{{ .Identifier }}" {
	name 			= "${terraform.workspace}-{{.Name}}"
	billing_mode 	= "PAY_PER_REQUEST"
	hash_key 		= "{{.HashKey.Field}}"
	{{ if .SortKey -}}
	range_key		= "{{.SortKey.Field}}"
	attribute {
		name = "{{.SortKey.Field}}"
		type = "{{.SortKey.Type}}"
	}
	{{- end }}
	attribute {
		name = "{{.HashKey.Field}}"
		type = "{{.HashKey.Type}}"
	}
	ttl {
		attribute_name = "" # Has to be empty or terraform won't update properly
		enabled        = false
	}
	tags = {
		Environment = terraform.workspace
	}
}
resource "aws_appsync_datasource" "{{ .Identifier }}" {
	api_id 				= aws_appsync_graphql_api.appsync.id
	name 				= "${terraform.workspace}_{{ .Name }}"
	service_role_arn 	= aws_iam_role.appsync.arn
	type				= "AMAZON_DYNAMODB"
	depends_on			= [
		aws_dynamodb_table.{{ .Identifier }}
	]
	dynamodb_config {
		table_name = aws_dynamodb_table.{{ .Identifier }}.name
	}
}
{{end}}
# ===============================================
# RESOLVERS
# ===============================================
{{range .Resolvers}}
resource "aws_appsync_resolver" "{IDENTIFIER}" {
	api_id

	# TODO
}
{{end}}
`))
