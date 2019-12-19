package graphql_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/graphql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalField(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected *graphql.Field
		err      error
	}{
		{
			"Single field with type",
			[]byte("name: id\ntype: ID"),
			&graphql.Field{
				Name:     "id",
				Type:     &graphql.FieldType{"ID", false, false},
				Resolver: nil,
			},
			nil,
		},
		{
			"Single non-nullable field with type",
			[]byte("name: id\ntype: ID!"),
			&graphql.Field{
				Name:     "id",
				Type:     &graphql.FieldType{"ID", false, true},
				Resolver: nil,
			},
			nil,
		},
		{
			"Field with default type",
			[]byte("name: default-type"),
			&graphql.Field{
				Name:     "default-type",
				Type:     &graphql.FieldType{"String", false, false},
				Resolver: nil,
			},
			nil,
		},
		{
			"Field with input type",
			[]byte("name: type-with-input\ninputType: Int"),
			&graphql.Field{
				Name:      "type-with-input",
				Type:      &graphql.FieldType{"String", false, false},
				Resolver:  nil,
				InputType: &graphql.FieldType{"Int", false, false},
			},
			nil,
		},
		{
			"Field with object type",
			[]byte("name: channel\ntype: Channel"),
			&graphql.Field{
				Name:     "channel",
				Type:     &graphql.FieldType{"Channel", false, false},
				Resolver: nil,
			},
			nil,
		},
		{
			"Field with list type",
			[]byte("name: list-type\ntype: [String]"),
			&graphql.Field{
				Name:     "list-type",
				Type:     &graphql.FieldType{"String", true, false},
				Resolver: nil,
			},
			nil,
		},
		{
			"Field with non-nullable list type",
			[]byte("name: list-type\ntype: [String!]"),
			&graphql.Field{
				Name:     "list-type",
				Type:     &graphql.FieldType{"String", true, true},
				Resolver: nil,
			},
			nil,
		},
		{
			"Field with resolver",
			[]byte("name: bad\nresolver:\n  action: get\n"),
			&graphql.Field{
				Name:     "bad",
				Type:     nil,
				Resolver: &graphql.Resolver{Action: graphql.ActionGet},
			},
			nil,
		},
		{
			"Field should error with type and resolver",
			[]byte("name: bad\ntype: String\nresolver:\n  action: get\n"),
			&graphql.Field{
				Name:     "bad",
				Type:     nil,
				Resolver: nil,
			},
			graphql.ErrTypeAndResolver,
		},
	} {
		var f graphql.Field
		err := yaml.Unmarshal(c.yaml, &f)
		switch c.err {
		case nil:
			assert.NoError(t, err)
			assert.Equal(t, c.expected, &f, c.scenario)
		default:
			assert.EqualError(t, err, c.err.Error())
		}

	}
}
