package manifest_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestResolverUnmarshalYAML(t *testing.T) {
	type testcase struct {
		scenario string
		yaml     []byte
		expected *manifest.Resolver
		err      error
	}

	cases := []testcase{
		{
			scenario: "Can't supply template with defined action",
			yaml:     []byte("source: AnimalTable\naction: get-item\ntemplate: my-template"),
			expected: nil,
			err:      errors.New("must not specify a custom template when resolver action is not 'custom'"),
		},
		{
			scenario: "Can't use custom without template (no key)",
			yaml:     []byte("source: AnimalTable\naction: custom"),
			expected: nil,
			err:      errors.New("must specify custom template when resolver action is 'custon'"),
		},
		{
			scenario: "Can't use custom without template (empty value)",
			yaml:     []byte("source: AnimalTable\naction: custom\ntemplate:"),
			expected: nil,
			err:      errors.New("must specify custom template when resolver action is 'custon'"),
		},
		{
			scenario: "Custom template OK",
			yaml:     []byte("source: AnimalTable\naction: custom\ntemplate: my-template"),
			expected: &manifest.Resolver{
				Source:   "AnimalTable",
				Action:   "custom",
				Template: "my-template",
			},
			err: nil,
		},
		{
			scenario: "Standard action OK",
			yaml:     []byte("source: AnimalTable\naction: get-item"),
			expected: &manifest.Resolver{
				Source:   "AnimalTable",
				Action:   "get-item",
				Template: "get-item",
			},
			err: nil,
		},
	}

	for _, c := range cases {
		var r manifest.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		switch c.err {
		case nil:
			assert.NoError(t, err)
			assert.Equal(t, c.expected, &r, c.scenario)
		default:
			assert.EqualError(t, err, c.err.Error(), c.scenario)
		}
	}
}
