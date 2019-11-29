package schema_test

import (
	"testing"

	"github.com/ONSdigital/appsync-resolver-builder/pkg/schema"
)

func TestNew(t *testing.T) {
	_, err := schema.New(definition)
	if err != nil {
		t.Errorf("error parsing definition '%v', expected nil", err)
	}
}

var definition []byte = []byte(`
enums:
  - name: Channel
    values: [EMAIL,LETTER,OTHER]

objects:
- name: Correspondence
  fields:
    - name: reference
      type: ID
      nonNullable: true
    - name: subject
    - name: enquiry
    - name: Channel
      type: Channel
`)
