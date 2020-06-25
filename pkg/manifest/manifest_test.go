package manifest_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestCheckVersion(t *testing.T) {

	tests := []struct {
		Data   []byte
		ErrStr string
	}{
		{[]byte("---\n"), "no valid 'version' found in manifest - must be in format 'v1.0.0'"},
		{[]byte("---\n1.0.0"), "no valid 'version' found in manifest - must be in format 'v1.0.0'"},
		{[]byte("---\nversion: v2.0.0\n"), ""},
		{[]byte("---\nversion: v2\n"), ""},
		{[]byte("---\nversion: v2.4.0\n"), ""},
		{[]byte("---\nversion: v2.0.99\n"), ""},
		{[]byte("---\nversion: v1.0.99\n"), "manifest must be version >= 2.0.0 and < 3.0.0, got v1.0.99"},
		{[]byte("---\nversion: v3.0.0\n"), "manifest must be version >= 2.0.0 and < 3.0.0, got v3.0.0"},
	}

	for _, test := range tests {
		err := manifest.CheckVersion(test.Data)
		if test.ErrStr != "" {
			assert.EqualError(t, err, test.ErrStr)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestFieldIsNonNullable(t *testing.T) {

	tests := []struct {
		Field       manifest.Field
		NonNullable bool
	}{
		{manifest.Field{Name: "defaultField"}, false},
		{manifest.Field{Name: "not:String"}, false},
		{manifest.Field{Name: "nonNullable:String!"}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.Field.IsNonNullable(), test.NonNullable)
	}
}

func TestDynamoDataSource(t *testing.T) {

	tests := []struct {
		Yaml   []byte
		ErrStr string
	}{
		{[]byte("---\nversion: v2.0.0\napiNameSuffix: something\nsources:\n  MySource:\n    dynamo:\n      hashKey: cat"), ""},
	}

	for _, test := range tests {
		_, err := manifest.Parse(test.Yaml)
		if test.ErrStr != "" {
			assert.EqualError(t, err, test.ErrStr)
		} else {
			assert.NoError(t, err)
		}
	}

}

func TestGettingAttributeNameAndType(t *testing.T) {

	// Tests cover:
	// - GetAttributeName
	// - GetAttributeType
	// - GetAttributeTypeStripped

	tests := []struct {
		scenario         string
		name             string
		expectedName     string
		expectedType     string
		expectedStripped string
	}{
		{
			scenario:         "Simple case with default type",
			name:             "name",
			expectedName:     "name",
			expectedType:     "String",
			expectedStripped: "String",
		},
		{
			scenario:         "Simple case with defined type",
			name:             "name:Boolean!",
			expectedName:     "name",
			expectedType:     "Boolean!",
			expectedStripped: "Boolean",
		},
		{
			scenario:         "List type",
			name:             "name:[String]",
			expectedName:     "name",
			expectedType:     "[String]",
			expectedStripped: "String",
		},
		{
			scenario:         "List mandatory type",
			name:             "name:[String!]",
			expectedName:     "name",
			expectedType:     "[String!]",
			expectedStripped: "String",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expectedName, manifest.GetAttributeName(test.name), test.scenario)
		assert.Equal(t, test.expectedType, manifest.GetAttributeType(test.name, "String"), test.scenario)
		assert.Equal(t, test.expectedStripped, manifest.GetAttributeTypeStripped(test.name, "String"), test.scenario)
	}
}
