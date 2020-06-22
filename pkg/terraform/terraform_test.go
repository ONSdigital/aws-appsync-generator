package terraform_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/terraform"
	"github.com/stretchr/testify/assert"
)

func TestNewDynamoKey(t *testing.T) {
	tests := []struct {
		In       string
		Expected *terraform.DynamoKey
	}{
		{"pk", &terraform.DynamoKey{"pk", "S"}}, // Test that default should be "S"
		{"pk:S", &terraform.DynamoKey{"pk", "S"}},
		{"number:N", &terraform.DynamoKey{"number", "N"}},
		{"", nil},
	}
	for _, test := range tests {
		dk := terraform.NewDynamoKey(test.In)
		if test.Expected != nil {
			assert.Equal(t, test.Expected.Field, dk.Field)
			assert.Equal(t, test.Expected.Type, dk.Type)
		} else {
			assert.Nil(t, dk)
		}

	}
}
