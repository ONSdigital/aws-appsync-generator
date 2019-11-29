package schema_test

import (
	"testing"

	"github.com/ONSdigital/appsync-resolver-builder/pkg/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestIsLegalScalar(t *testing.T) {
	types := []string{
		"String", "Int", "ID", "Float", "Boolean",
		"AWSDate", "AWSTime", "AWSDateTime", "AWSTimestamp",
		"AWSEmail", "AWSJSON", "AWSURL", "AWSPhone", "AWSIPAddress",
	}
	for _, tp := range types {
		f := &schema.Field{Type: tp}
		assert.Truef(t, f.IsLegalScalarType(), "Type '%s' should be legal type", tp)
	}

	f := schema.Field{Type: "Complex"}
	assert.False(t, f.IsLegalScalarType(), "Type 'Complex' should not be a legal type")
}

func TestUnmarshalField(t *testing.T) {
	for _, c := range []struct {
		yaml     []byte
		expected *schema.Field
	}{
		{
			[]byte("name: id\ntype: ID\nnonNullable: true"),
			&schema.Field{
				Name:        "id",
				Type:        "ID",
				NonNullable: true,
				Resolver:    nil,
			},
		},
		{
			[]byte("name: default-type"),
			&schema.Field{
				Name:        "default-type",
				Type:        "String",
				NonNullable: false,
				Resolver:    nil,
			},
		},
		{
			[]byte("name: type-with-input\ninputType: Int"),
			&schema.Field{
				Name:        "type-with-input",
				Type:        "String",
				NonNullable: false,
				Resolver:    nil,
				InputType:   "Int",
			},
		},
	} {
		var f schema.Field
		err := yaml.Unmarshal(c.yaml, &f)
		assert.NoErrorf(t, err, "expected no error, got '%v'", err)
		assert.Equal(t, c.expected, &f, "unmarshal incorrect")
	}
}
