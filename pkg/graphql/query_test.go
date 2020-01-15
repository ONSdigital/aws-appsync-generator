package graphql_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/graphql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalQuery(t *testing.T) {
	for _, c := range []struct {
		yaml     []byte
		expected *graphql.Field
	}{
		{
			[]byte("name: correspondence\ntype: Correspondence\nresolver:\n  action: get\n  keyFields:\n    - name: reference"),
			&graphql.Field{
				Name: "correspondence",
				// Resolver: &graphql.GetResolver{
				// 	Action: "get",
				// },
				// Type:        &graphql.Query{"ID", false},
				// NonNullable: true,
				// Resolver:    nil,
			},
		},
	} {
		var q graphql.Query
		err := yaml.Unmarshal(c.yaml, &q)
		assert.NoErrorf(t, err, "expected no error, got '%v'", err)
		// assert.IsType(t, graphql.GetResolver{}, q.Resolver)
		// assert.Equal(t, c.expected, &f, "unmarshal incorrect")
	}
}

func TestNewInputFromObject(t *testing.T) {

	object := graphql.Object{
		Name: "Cat",
		Fields: []*graphql.Field{
			{
				Name: "fur",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      false,
					NonNullable: false,
				},
				InputType: &graphql.FieldType{
					Name:        "Int",
					IsList:      false,
					NonNullable: false,
				},
			},
			{
				Name: "name",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      false,
					NonNullable: true,
				},
			},
			{
				Name: "toys",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      true,
					NonNullable: false,
				},
			},
		},
	}

	expected := &graphql.InputObject{
		Name: "CreateCatInput",
		Fields: []*graphql.Field{
			{
				Name: "fur",
				Type: &graphql.FieldType{
					Name:        "Int",
					IsList:      false,
					NonNullable: false,
				},
			},
			{
				Name: "name",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      false,
					NonNullable: false,
				},
			},
			{
				Name: "toys",
				Type: &graphql.FieldType{
					Name:        "String",
					IsList:      true,
					NonNullable: false,
				},
			},
		},
	}

	input, err := graphql.NewInputFromObject(&object, graphql.ActionInsert)

	assert.NoError(t, err)
	assert.Equal(t, expected, input)

}
