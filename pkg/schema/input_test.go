package schema_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalInput(t *testing.T) {
	for _, c := range []struct {
		yaml     []byte
		expected *schema.Input
		errMsg   string
	}{
		{
			[]byte("name: MyInputWithBase\nbase: MyObject"),
			&schema.Input{
				Name:              "MyInputWithBase",
				Base:              "MyObject",
				NonNullableParams: nil,
				Exclude:           nil,
				Params:            nil,
			},
			"",
		},
		{
			[]byte("name: MyInputWithParams\nparams:\n  - name: a"),
			&schema.Input{
				Name:              "MyInputWithParams",
				Base:              "",
				NonNullableParams: nil,
				Exclude:           nil,
				Params: []*schema.ResolverParam{
					{
						Name:       "a",
						Type:       "String",
						SourceName: "",
						Fuzzy:      false,
					},
				},
			},
			"",
		},
		{
			[]byte("name: MyInput\nbase: MyObject\nparams:\n  - name: someField\n"),
			nil,
			"can only specify one of 'base' or 'params' in input 'MyInput'",
		},
	} {
		var i schema.Input
		err := yaml.Unmarshal(c.yaml, &i)
		if c.errMsg == "" {
			// Expecting no marshal error
			assert.NoErrorf(t, err, "expected no error, got '%v'", err)
			assert.Equal(t, c.expected, &i, "not unmarshaled as expected")
		} else {
			// Expecing a marshaling error
			assert.EqualError(t, err, c.errMsg, "expected error '%v', got '%v'", c.errMsg, err)
		}
	}
}
