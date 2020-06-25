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
				DataSourceName: "AnimalTable",
				Action:         "custom",
				Template:       "my-template",
			},
			err: nil,
		},
		{
			scenario: "Standard action OK",
			yaml:     []byte("source: AnimalTable\naction: get-item"),
			expected: &manifest.Resolver{
				DataSourceName: "AnimalTable",
				Action:         "get-item",
				Template:       "get-item",
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

func TestArgsSource(t *testing.T) {
	type tc struct {
		scenario string
		resolver manifest.Resolver
		expected string
	}

	cases := []tc{
		{
			scenario: "Query should be 'args'",
			resolver: manifest.Resolver{ParentType: "Query"},
			expected: "args",
		},
		{
			scenario: "Mutation should be 'args'",
			resolver: manifest.Resolver{ParentType: "Mutation"},
			expected: "args",
		},
		{
			scenario: "Object should be 'source'",
			resolver: manifest.Resolver{ParentType: "Animal"},
			expected: "source",
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, c.resolver.ArgsSource(), c.scenario)
	}
}

func TestKeyFieldJSONMapAndList(t *testing.T) {

	// Using the same scenarios to test both
	// - KeyFieldJSONMap()
	// - KeyFieldJSONList()

	type tc struct {
		scenario     string
		resolver     manifest.Resolver
		expectedMap  string
		expectedList string
	}

	cases := []tc{
		{
			scenario:     "No fields",
			resolver:     manifest.Resolver{},
			expectedMap:  "{}",
			expectedList: "[]",
		},
		{
			scenario: "One mandatory field",
			resolver: manifest.Resolver{
				KeyFields: []manifest.Field{
					{Name: "id:ID!"},
				},
			},
			expectedMap:  `{"id":"ID"}`,
			expectedList: `["id"]`,
		},
		{
			scenario: "Several fields (with defaults)",
			resolver: manifest.Resolver{
				KeyFields: []manifest.Field{
					{Name: "id:ID!"},
					{Name: "color"},
					{Name: "size:Int"},
				},
			},
			expectedMap:  `{"id":"ID","color":"String","size":"Int"}`,
			expectedList: `["id","color","size"]`,
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expectedMap, c.resolver.KeyFieldJSONMap(), c.scenario+" (map)")
		assert.Equal(t, c.expectedList, c.resolver.KeyFieldJSONList(), c.scenario+" (list)")
	}
}
