package schema_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/stretchr/testify/assert"
)

func TestStringifyQuery(t *testing.T) {
	type tc struct {
		scenario string
		query    *schema.Query
		expected string
	}

	cases := []tc{
		{
			scenario: "Simple with signature",
			query:    &schema.Query{Name: "animal", ReturnType: "Animal", Signature: "id: ID!"},
			expected: "animal(id: ID!): Animal",
		},
		{
			scenario: "Simple with no signature",
			query:    &schema.Query{Name: "animal", ReturnType: "Animal", Signature: ""},
			expected: "animal: Animal",
		},
	}

	for _, c := range cases {
		str := c.query.String()
		assert.Equal(t, c.expected, str, c.scenario)
	}
}
