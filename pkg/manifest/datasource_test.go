package manifest_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/stretchr/testify/assert"
)

func TestValidateDataSource(t *testing.T) {
	missingHashKey := "dynamo datasource must declare a hashKey"
	tests := []struct {
		DataSource manifest.DataSource
		ErrStr     string
	}{
		{manifest.DataSource{Dynamo: &manifest.DynamoSource{HashKey: "id:S"}}, ""},
		{manifest.DataSource{Dynamo: &manifest.DynamoSource{}}, missingHashKey},
		{manifest.DataSource{Dynamo: &manifest.DynamoSource{SortKey: "sorted"}}, missingHashKey},
	}
	for _, test := range tests {
		err := test.DataSource.Validate()
		if test.ErrStr == "" {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, test.ErrStr, err.Error())
		}
	}
}
