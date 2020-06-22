package graphql_test

import (
	"errors"
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg_old/graphql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalDataSource(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected *graphql.Source
		err      error
	}{
		{
			"Good dynamo data source",
			[]byte("name: dynamosource\ndynamo:\n  hash_key:\n    name: key"),
			&graphql.Source{
				Name: "dynamosource",
				Type: "dynamo",
				Dynamo: &graphql.DynamoSource{
					HashKey: &graphql.DynamoKeyType{
						Name: "key",
						Type: "S",
					},
				},
			},
			nil,
		},
		{
			"Good aurora data source",
			[]byte("name: aurorasource\nsql:\n  primary_key: key"),
			&graphql.Source{
				Name: "aurorasource",
				Type: "sql",
				SQL: &graphql.SQLSource{
					PrimaryKey: "key",
				},
			},
			nil,
		},
		// {
		// 	"Unsupported type",
		// 	[]byte("name: unsupported\ntype: sheepdb"),
		// 	nil,
		// 	errors.New("datasource 'unsupported' has unsupported type: sheepdb"),
		// },
		{
			"Missing name",
			[]byte("type: dynamo"),
			nil,
			errors.New("datasource does not declare a name"),
		},
	} {
		var s graphql.Source
		err := yaml.Unmarshal(c.yaml, &s)
		switch c.err {
		case nil:
			assert.NoError(t, err)
			assert.Equal(t, c.expected, &s, c.scenario)
		default:
			assert.EqualError(t, err, c.err.Error())
		}
	}
}
