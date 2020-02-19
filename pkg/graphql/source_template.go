package graphql

var sourceTemplate = `
{{ if eq .Type "dynamo" -}}
resource "aws_iam_role_policy" "record_dynamo_{{.Name}}" {
	name		= "${terraform.workspace}-dynamo-{{.Name}}"
	role 		= aws_iam_role.record.id
	policy 		= <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
    "Action": [
      "dynamodb:*"
    ],
    "Effect": "Allow",
    "Resource": [
      "${aws_dynamodb_table.{{.Name}}.arn}"
    ]
    }
  ]
}
EOF
  }

resource "aws_dynamodb_table" "{{.Name}}" {
	name 			= "${terraform.workspace}-{{.Name}}"
	billing_mode 	= "PAY_PER_REQUEST"
	hash_key 		= "{{.Dynamo.HashKey.Name}}"
	{{ if .Dynamo.SortKey -}}
	range_key		= "{{.Dynamo.SortKey.Name}}"

	attribute {
		name = "{{.Dynamo.SortKey.Name}}"
		type = "{{.Dynamo.SortKey.Type}}"
	}
	{{- end }}

	attribute {
		name = "{{.Dynamo.HashKey.Name}}"
		type = "{{.Dynamo.HashKey.Type}}"
	}

	ttl {
		attribute_name = "" # Has to be empty or terraform won't update properly
		enabled        = false
	}

	tags = {
		Environment = terraform.workspace
		Name        = "{{.Name}}"
	}
}

resource "aws_appsync_datasource" "{{.Name}}" {
	api_id 				= aws_appsync_graphql_api.record.id
	name 				= "${terraform.workspace}_{{.Name}}"
	service_role_arn 	= aws_iam_role.record.arn
	type				= "AMAZON_DYNAMODB"
	depends_on			= [
		aws_dynamodb_table.{{.Name}}
	]
	dynamodb_config {
		table_name = aws_dynamodb_table.{{.Name}}.name
	}
}
{{- end }}
`
