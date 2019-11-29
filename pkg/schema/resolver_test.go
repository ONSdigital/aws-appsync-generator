package schema_test

import (
	"testing"

	"github.com/ONSdigital/appsync-resolver-builder/pkg/schema"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestUnmarshalResolver(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected schema.Resolver
	}{
		{
			"Standard decode",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, type: String}, {name: cute, type: String}]}`),
			schema.Resolver{
				Action: "get",
				Params: []*schema.ResolverParam{
					{Name: "animal", Type: "String", NonNullable: true},
					{Name: "cute", Type: "String", NonNullable: false}},
			},
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")
		assert.Equal(t, c.expected, r, "incorrect unmarshal")
	}
}

func TestAppsyncFieldsMapString(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected string
	}{
		{
			"Standard decode",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, type: String}, {name: cute, type: String}]}`),
			`{"animal": "animal","cute": "cute"}`,
		},
		{
			"Default param type",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true}, {name: cute, type: String}]}`),
			`{"animal": "animal","cute": "cute"}`,
		},
		{
			"With Source name",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, sourceName: sheep}, {name: cute, type: String}]}`),
			`{"animal": "sheep","cute": "cute"}`,
		},
		{
			"No params",
			[]byte(`{action: get}`),
			"{}",
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")
		assert.Equal(t, c.expected, r.AppsyncFieldsMapString(), "correct string")
	}
}

func TestAppsyncFuzzyMapString(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected string
	}{
		{
			"Standard decode",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, type: String, fuzzy: true}, {name: cute, type: String}]}`),
			`{"animal": true,"cute": false}`,
		},
		{
			"No params",
			[]byte(`{action: get}`),
			"{}",
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")
		assert.Equal(t, c.expected, r.AppsyncFuzzyMapString(), "correct string")
	}
}

func TestTemplateNames(t *testing.T) {
	for _, c := range []struct {
		scenario         string
		yaml             []byte
		expectedError    string
		expectedRequest  string
		expectedResponse string
	}{
		{
			"Unknown action",
			[]byte(`{action: immolate}`),
			"action type 'immolate' not allowable for schema.Resolver",
			"",
			"",
		},
		{
			"Standard action",
			[]byte(`{action: get}`),
			"",
			"get.tmpl",
			"get.tmpl",
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")

		rq, rs, err := r.TemplateNames()
		if c.expectedError == "" {
			assert.Equal(t, rq, c.expectedRequest, "request template")
			assert.Equal(t, rs, c.expectedResponse, "response template")
			assert.NoError(t, err, "no error")
		} else {
			assert.Empty(t, rq, "no request template")
			assert.Empty(t, rs, "no response template")
			assert.EqualError(t, err, "action type 'immolate' not allowable for resolver", "expected error")
		}
	}
}

func TestParamString(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected string
	}{
		{
			"Standard decode",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, type: String}, {name: cute, type: String}]}`),
			"animal: String!, cute: String",
		},
		{
			"Default param type",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true}, {name: cute, type: String}]}`),
			"animal: String!, cute: String",
		},
		{
			"No params",
			[]byte(`{action: get}`),
			"",
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")
		assert.Equal(t, c.expected, r.ParamString(), "correct string")
	}
}

func TestFieldsString(t *testing.T) {
	for _, c := range []struct {
		scenario string
		yaml     []byte
		expected string
	}{
		{
			"Standard decode",
			[]byte(`{action: get, params: [{name: animal, nonNullable: true, type: String}, {name: cute, type: String}]}`),
			"animal,cute",
		},
		{
			"No fields",
			[]byte(`{action: get}`),
			"",
		},
	} {
		var r schema.Resolver
		err := yaml.Unmarshal(c.yaml, &r)
		assert.NoError(t, err, "got error marshaling, expected nil")
		assert.Equal(t, c.expected, r.FieldsString(), "correct string")
	}
}
