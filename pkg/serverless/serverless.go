package serverless

import (
	"fmt"
	"io"
	"strings"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Constants defining the valid types for data sources
const (
	DynamoSourceType        = "AMAZON_DYNAMODB"
	RDSSourceType           = "RELATIONAL_DATABASE"
	LambdaSourceType        = "AWS_LAMBDA"
	ElasticSearchSourceType = "AMAZON_ELASTICSEARCH"
	HTTPSourceType          = "HTTP"
)

type (
	// Configuration represents a serverless configuration yaml
	Configuration struct {
		Provider  ProviderBlock `yaml:"provider"`
		Service   string        `yaml:"service"`
		Plugins   []string      `yaml:"plugins"`
		Custom    CustomBlock   `yaml:"custom"`
		Resources ResourceBlock `yaml:"resources"`
	}

	ProviderBlock struct {
		Name   string `yaml:"name"`
		Stage  string `yaml:"stage"`
		Region string `yaml:"region"`
	}

	ResourceBlock struct {
		Resources map[string]Resource `yaml:"Resources"`
	}

	LogConfig struct {
		CloudWatchLogsRoleArn FnGettAtt `yaml:"CloudWatchLogsRoleArn"`
		Level                 string    `yaml:"level"`
	}

	MappingTemplate struct {
	}

	CustomBlock struct {
		AppSync AppSyncBlock `yaml:"appSync"`
	}

	AppSyncBlock struct {
		Name               string            `yaml:"name,omitempty"`
		Schema             string            `yaml:"schema"`
		AuthenticationType string            `yaml:"authenticationType"`
		LogConfig          LogConfig         `yaml:"logConfig"`
		MappingTemplates   []MappingTemplate `yaml:"mappingTemplates"`
		DataSources        []DataSource      `yaml:"dataSources"`
	}

	// ==== Data Sources ========================

	// DataSource represents an appsync datasource
	DataSource struct {
		Type        string      `yaml:"type"`
		Name        string      `yaml:"name"`
		Description string      `yaml:"description,omitempty"`
		Config      interface{} `yaml:"config"`
	}

	// DataSourceConfig is the common elements of data source configurations
	DataSourceConfig struct {
		IAMRoleStatements []IAMPolicyDocumentStatement `yaml:"iamRoleStatements,omitempty"`
		Region            string                       `yaml:"region,omitempty"`
	}

	// DynamoDataSourceConfig represents the configuration of a dynamo datasource
	DynamoDataSourceConfig struct {
		DataSourceConfig `yaml:",inline"`
		TableName        string    `yaml:"tableName"`
		ServiceRoleArn   FnGettAtt `yaml:"serviceRoleArn,omitempty"`
	}

	// RDSDataSourceConfig represents the configuration for an RDS datasource
	RDSDataSourceConfig struct {
		DataSourceConfig  `yaml:",inline"`
		ClusterIdentifier CFFunction `yaml:"dbClusterIdentifier"`
		AWSSecretStoreARN CFFunction `yaml:"awsSecretStoreArn"`
		ServiceRoleArn    FnGettAtt  `yaml:"serviceRoleArn"`
		DatabaseName      string     `yaml:"databaseName,omitempty"`
		Schema            string     `yaml:"schema,omitempty"`
	}

	// ==== Resources ===========================

	// Resource is a CloudFormation resource
	Resource struct {
		Type       string             `yaml:"Type"`
		Properties ResourceProperties `yaml:"Properties"`
	}

	// ResourceProperties are CloudFormation resources
	ResourceProperties interface{}

	// ==== DyanmoDB Resources ==================

	DynamoResourceProperties struct {
		TableName            string                         `yaml:"TableName"`
		AttributeDefinitions []DynamoAttributeDefinition    `yaml:"AttributeDefinitions"`
		KeySchema            []DynamoKeyAttributeDefinition `yaml:"KeySchema"`
		BillingMode          string                         `yaml:"BillingMode"`
	}

	DynamoAttributeDefinition struct {
		AttributeName string `yaml:"AttributeName"`
		AttributeType string `yaml:"AttributeType"`
	}

	DynamoKeyAttributeDefinition struct {
		AttributeName string `yaml:"AttributeName"`
		KeyType       string `yaml:"KeyType"`
	}

	// ==== IAM Resources =======================

	// IAMResourceProperties are CloudFormation IAM resource properties
	IAMResourceProperties struct {
		Path                     string              `yaml:"Path"`
		RoleName                 string              `yaml:"RoleName"`
		AssumeRolePolicyDocument IAMPolicyDocument   `yaml:"AssumeRolePolicyDocument"`
		ManagedPolicyArns        []string            `yaml:"ManagedPolicyArns,omitempty"`
		Policies                 []IAMResourcePolicy `yaml:"Policies,omitempty"`
	}

	// IAMResourcePolicy is an IAM policy
	IAMResourcePolicy struct {
		PolicyName     string            `yaml:"PolicyName"`
		PolicyDocument IAMPolicyDocument `yaml:"PolicyDocument"`
	}

	// IAMPolicyDocument is an IAM policy document
	IAMPolicyDocument struct {
		Version   string                       `yaml:"Version"`
		Statement []IAMPolicyDocumentStatement `yaml:"Statement"`
	}

	// IAMPolicyDocumentStatement is an IAM policy document statement
	IAMPolicyDocumentStatement struct {
		Effect    string                              `yaml:"Effect"`
		Action    []string                            `yaml:"Action"`
		Resource  []string                            `yaml:"Resource,omitempty"`
		Principal IAMPolicyDocumentStatementPrincipal `yaml:"Principal,omitempty"`
	}

	// IAMPolicyDocumentStatementPrincipal is a principal service in a policy document
	IAMPolicyDocumentStatementPrincipal struct {
		Service []string `yaml:"Service"`
	}

	// --- Functions

	CFFunction interface{}

	// FnGettAtt is a Fn::GetAtt function call
	// { Fn::GetAtt: [ITEM, Arn] }
	FnGettAtt struct {
		Params []string `yaml:"Fn::GetAtt"`
	}

	// FnRef is a Ref function call
	// { Ref: ITEM }
	FnRef struct {
		Ref string `yaml:"Ref"`
	}
)

// New returns a new empty Configuration
func New() *Configuration {
	return &Configuration{
		Provider: ProviderBlock{
			Name:   "aws",
			Stage:  "${opt:stage, 'dev'}",
			Region: "${env:AWS_REGION}",
		},
		Plugins: []string{"serverless-appsync-plugin"},
		Custom: CustomBlock{
			AppSyncBlock{
				Name:               "",
				Schema:             "schema.public.graphql",
				AuthenticationType: "API_KEY",
				LogConfig: LogConfig{
					CloudWatchLogsRoleArn: FnGettAtt{
						Params: []string{"AppSyncLoggingServiceRole", "Arn"},
					},
					Level: "ERROR",
				},
			},
		},
		Resources: ResourceBlock{
			Resources: make(map[string]Resource),
		},
	}
}

// NewFromManifest returns a populated Configuration build from the
// given manifest data
func NewFromManifest(m *manifest.Manifest) (*Configuration, error) {

	sc := New()
	sc.Custom.AppSync.Name = m.APINameSuffix
	sc.Service = "${env:DEPLOYMENT}-" + stripNameToIdentifier(m.APINameSuffix)

	// Default IAM Policies
	sc.Resources.Resources["AppSyncLoggingServiceRole"] = Resource{
		Type: "AWS::IAM::Role",
		Properties: IAMResourceProperties{
			Path:     "/appsync/log/role/",        // TODO what should this really be?
			RoleName: "AppSyncLoggingServiceRole", // TODO needs to namespaced?
			AssumeRolePolicyDocument: IAMPolicyDocument{
				Version: "2012-10-17",
				Statement: []IAMPolicyDocumentStatement{
					{
						Effect: "Allow",
						Action: []string{"sts:AssumeRole"},
						Principal: IAMPolicyDocumentStatementPrincipal{
							Service: []string{"appsync.amazonaws.com"},
						},
					},
				},
			},
			ManagedPolicyArns: []string{
				"arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs",
			},
		},
	}
	// sc.Resources.Resources["AppSyncDynamoDBServiceRole"] = Resource{
	// 	Type: "AWS::IAM::Role",
	// 	Properties: IAMResourceProperties{
	// 		Path:     "/appsync/log/role/",         // TODO what should this really be?
	// 		RoleName: "AppSyncDynamoDBServiceRole", // TODO needs to namespaced?
	// 		AssumeRolePolicyDocument: IAMPolicyDocument{
	// 			Version: "2012-10-17",
	// 			Statement: []IAMPolicyDocumentStatement{
	// 				{
	// 					Effect: "Allow",
	// 					Action: []string{"sts:AssumeRole"},
	// 					Principal: IAMPolicyDocumentStatementPrincipal{
	// 						Service: []string{"appsync.amazonaws.com"},
	// 					},
	// 				},
	// 			},
	// 		},
	// 		Policies: []IAMResourcePolicy{
	// 			{
	// 				PolicyName: m.APINameSuffix + "-dynamodb",
	// 				PolicyDocument: IAMPolicyDocument{
	// 					Version: "2012-10-17",
	// 					Statement: []IAMPolicyDocumentStatement{
	// 						{
	// 							Effect:   "Allow",
	// 							Action:   []string{"dynamodb:*"},
	// 							Resource: []string{"arn:aws:dynamodb:::table/" + m.APINameSuffix + "-*"},
	// 							// Resource: []string{"*"},
	// 							// TODO:: Property validation failure: [Value for property {/ServiceRoleArn} does not match pattern {^arn:aws(-cn)?:iam::\d{12}:role\/[a-zA-Z0-9=,_.@/-]{1,2048}
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// Import data sources
	for dsName, ds := range m.DataSources {
		source := ds.GetSource()
		switch t := source.(type) {
		case *manifest.DynamoSource:

			resourceIdentifier := "DynamoTable" + dsName
			tableName := m.APINameSuffix + "-" + dsName

			// Add the datasource reference to Appsync
			sc.Custom.AppSync.DataSources = append(sc.Custom.AppSync.DataSources, DataSource{
				Type: DynamoSourceType,
				Name: dsName,
				Config: DynamoDataSourceConfig{
					TableName: tableName,
					DataSourceConfig: DataSourceConfig{
						Region: "${env:AWS_REGION}",
						IAMRoleStatements: []IAMPolicyDocumentStatement{
							{
								Effect:   "Allow",
								Action:   []string{"dynamodb:*"},
								Resource: []string{"arn:aws:dynamodb:::table/" + tableName},
							},
						},
					},
				},
			})

			// Create the actual data source itself (TODO unless it's marked
			// as DO NOT CREATE)
			prop := DynamoResourceProperties{
				TableName: tableName,
				AttributeDefinitions: []DynamoAttributeDefinition{
					{
						AttributeName: manifest.GetAttributeName(ds.Dynamo.HashKey),
						AttributeType: manifest.GetAttributeType(ds.Dynamo.HashKey, "S"),
					},
				},
				KeySchema: []DynamoKeyAttributeDefinition{
					{
						AttributeName: manifest.GetAttributeName(ds.Dynamo.HashKey),
						KeyType:       "HASH",
					},
				},
				BillingMode: "PAY_PER_REQUEST",
			}

			if ds.Dynamo.SortKey != "" {
				prop.AttributeDefinitions = append(prop.AttributeDefinitions, DynamoAttributeDefinition{
					AttributeName: manifest.GetAttributeName(ds.Dynamo.SortKey),
					AttributeType: manifest.GetAttributeType(ds.Dynamo.SortKey, "S"),
				})
				prop.KeySchema = append(prop.KeySchema, DynamoKeyAttributeDefinition{
					AttributeName: manifest.GetAttributeName(ds.Dynamo.SortKey),
					KeyType:       "RANGE",
				})
			}

			sc.Resources.Resources[resourceIdentifier] = Resource{
				Type:       "AWS::DynamoDB::Table",
				Properties: prop,
			}

		case *manifest.RDSSource:

			resourceIdentifier := "RDSCluster" + dsName

			sc.Custom.AppSync.DataSources = append(sc.Custom.AppSync.DataSources, DataSource{
				Type: RDSSourceType,
				Name: dsName,
				Config: RDSDataSourceConfig{
					ClusterIdentifier: FnRef{resourceIdentifier},
					AWSSecretStoreARN: resourceIdentifier + "Secret",
					ServiceRoleArn:    FnGettAtt{},
				},
			})

		case *manifest.LambdaSource:

		default:
			return nil, fmt.Errorf("invalid data source type: %T", t)
		}
	}

	return sc, nil
}

// Marshal performs a YAML marshal on the configuration
func (sc *Configuration) Marshal() ([]byte, error) {
	return yaml.Marshal(sc)
}

// Write outputs the configuration to the io.Writer
func (sc *Configuration) Write(w io.Writer) error {
	data, err := sc.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to write serverless configuration")
	}
	_, err = w.Write(data)
	return err
}

func stripNameToIdentifier(name string) string {
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return strings.ToLower(name)
}
