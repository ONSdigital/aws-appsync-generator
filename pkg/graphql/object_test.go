package graphql_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/graphql"
	"github.com/stretchr/testify/assert"
)

func TestNewFilterFromObject(t *testing.T) {

	in := &graphql.Object{
		Name: "TestObject",
		Fields: []*graphql.Field{
			{
				Name: "name",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      false,
					NonNullable: true,
				},
			},
			{
				Name: "age",
				Type: &graphql.FieldType{
					Name:        "Int",
					IsList:      false,
					NonNullable: false,
				},
			},
			{
				Name: "custom",
				Type: &graphql.FieldType{
					Name:        "BananaObject",
					IsList:      false,
					NonNullable: false,
				},
			},
		},
	}

	expected := &graphql.FilterObject{
		Name: "TestObjectFilter",
		Fields: []*graphql.Field{
			{
				Name: "name",
				Type: &graphql.FieldType{
					Name:        "TableStringFilterInput",
					IsList:      false,
					NonNullable: false,
				},
			},
			{
				Name: "age",
				Type: &graphql.FieldType{
					Name:        "TableIntFilterInput",
					IsList:      false,
					NonNullable: false,
				},
			},
		},
	}

	generated := graphql.NewFilterFromObject(in)
	assert.Equal(t, expected, generated)
}
