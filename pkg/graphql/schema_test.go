package graphql_test

import (
	"testing"

	"github.com/ONSdigital/aws-appsync-generator/pkg/graphql"
	"github.com/stretchr/testify/assert"
)

func mustCompileSchema(t *testing.T, manifest []byte) *graphql.Schema {
	s, err := graphql.NewSchemaFromManifest(exampleSchemaManifest, "dyanmo")
	if err != nil {
		t.Fatalf("unable to parse manifest: %v", err)
	}
	return s
}

func TestNewSchemaFromManifest(t *testing.T) {
	_, err := graphql.NewSchemaFromManifest(exampleSchemaManifest, "dyanmo")
	if err != nil {
		t.Errorf("error parsing definition '%v', expected nil", err)
	}
}

func TestSchemaGenerateBytes(t *testing.T) {
	s := mustCompileSchema(t, exampleSchemaManifest)
	g, err := s.GenerateBytes()
	if err != nil {
		t.Errorf("error generating byte buffer '%v', expected nil", err)
	}
	assert.IsType(t, []byte{}, g)
	// spew.Dump(g)
}

var exampleSchemaManifest = []byte(`
enums:
  - name: Channel
    values: [EMAIL,LETTER,OTHER]

objects:
  - name: Correspondence
    fields:
      - name: reference
        type: ID!
      - name: subject
      - name: enquiry
      - name: channel
        type: Channel
      - name: copiedTo
        type: [String]
`)
