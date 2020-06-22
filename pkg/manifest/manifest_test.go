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
