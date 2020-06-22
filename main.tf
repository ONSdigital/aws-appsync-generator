
# ===============================================
# APPSYNC API
# ===============================================
resource "aws_iam_role" "appsync" {
	name = "${terraform.workspace}-testing_api"
	assume_role_policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Service\":\"appsync.amazonaws.com\"},\"Effect\":\"Allow\"}]}"
}
resource "aws_iam_role_policy_attachment" "appsync" {
	policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs"
	role       = aws_iam_role.appsync.name
}
resource "aws_appsync_graphql_api" "appsync" {
	authentication_type = "API_KEY"
	name                = "${terraform.workspace}-testing_api"
	schema              = "${file("schema.public.graphql")}"
	log_config {
		cloudwatch_logs_role_arn = aws_iam_role.appsync.arn
		field_log_level          = "ERROR"
	}
	tags = {
		Environment = terraform.workspace
		Deployment = "${terraform.workspace}-testing_api"
	}
}
output "graphql_api_id" {
	value = aws_appsync_graphql_api.appsync.id
}
output "graphql_host" {
	value = aws_appsync_graphql_api.appsync.uris
}

# ===============================================
# DATA SOURCES
# ===============================================

# Dynamo: SimpleTable ----
resource "aws_iam_role_policy" "record_dynamo_simpletable" {
	name		= "${terraform.workspace}-dynamo-SimpleTable"
	role 		= aws_iam_role.appsync.id
	policy 		= "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"dynamodb:*\"],\"Effect\":\"Allow\",\"Resource\":[\"${aws_dynamodb_table.simpletable.arn}\"]}]}"
}
resource "aws_dynamodb_table" "simpletable" {
	name 			= "${terraform.workspace}-SimpleTable"
	billing_mode 	= "PAY_PER_REQUEST"
	hash_key 		= "id"
	
	attribute {
		name = "id"
		type = "S"
	}
	ttl {
		attribute_name = "" # Has to be empty or terraform won't update properly
		enabled        = false
	}
	tags = {
		Environment = terraform.workspace
	}
}
resource "aws_appsync_datasource" "simpletable" {
	api_id 				= aws_appsync_graphql_api.appsync.id
	name 				= "${terraform.workspace}_SimpleTable"
	service_role_arn 	= aws_iam_role.appsync.arn
	type				= "AMAZON_DYNAMODB"
	depends_on			= [
		aws_dynamodb_table.simpletable
	]
	dynamodb_config {
		table_name = aws_dynamodb_table.simpletable.name
	}
}

# Dynamo: ActionNotes ----
resource "aws_iam_role_policy" "record_dynamo_actionnotes" {
	name		= "${terraform.workspace}-dynamo-ActionNotes"
	role 		= aws_iam_role.appsync.id
	policy 		= "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"dynamodb:*\"],\"Effect\":\"Allow\",\"Resource\":[\"${aws_dynamodb_table.actionnotes.arn}\"]}]}"
}
resource "aws_dynamodb_table" "actionnotes" {
	name 			= "${terraform.workspace}-ActionNotes"
	billing_mode 	= "PAY_PER_REQUEST"
	hash_key 		= "request_number"
	range_key		= "id"
	attribute {
		name = "id"
		type = "S"
	}
	attribute {
		name = "request_number"
		type = "S"
	}
	ttl {
		attribute_name = "" # Has to be empty or terraform won't update properly
		enabled        = false
	}
	tags = {
		Environment = terraform.workspace
	}
}
resource "aws_appsync_datasource" "actionnotes" {
	api_id 				= aws_appsync_graphql_api.appsync.id
	name 				= "${terraform.workspace}_ActionNotes"
	service_role_arn 	= aws_iam_role.appsync.arn
	type				= "AMAZON_DYNAMODB"
	depends_on			= [
		aws_dynamodb_table.actionnotes
	]
	dynamodb_config {
		table_name = aws_dynamodb_table.actionnotes.name
	}
}

# Dynamo: RefData ----
resource "aws_iam_role_policy" "record_dynamo_refdata" {
	name		= "${terraform.workspace}-dynamo-RefData"
	role 		= aws_iam_role.appsync.id
	policy 		= "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"dynamodb:*\"],\"Effect\":\"Allow\",\"Resource\":[\"${aws_dynamodb_table.refdata.arn}\"]}]}"
}
resource "aws_dynamodb_table" "refdata" {
	name 			= "${terraform.workspace}-RefData"
	billing_mode 	= "PAY_PER_REQUEST"
	hash_key 		= "refcat"
	range_key		= "order"
	attribute {
		name = "order"
		type = "N"
	}
	attribute {
		name = "refcat"
		type = "S"
	}
	ttl {
		attribute_name = "" # Has to be empty or terraform won't update properly
		enabled        = false
	}
	tags = {
		Environment = terraform.workspace
	}
}
resource "aws_appsync_datasource" "refdata" {
	api_id 				= aws_appsync_graphql_api.appsync.id
	name 				= "${terraform.workspace}_RefData"
	service_role_arn 	= aws_iam_role.appsync.arn
	type				= "AMAZON_DYNAMODB"
	depends_on			= [
		aws_dynamodb_table.refdata
	]
	dynamodb_config {
		table_name = aws_dynamodb_table.refdata.name
	}
}

# ===============================================
# RESOLVERS
# ===============================================

