package terraform

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
)

type (
	// Terraform represents the overall data to populate the terraform
	Terraform struct {
		AppSync     AppSync
		Resolvers   []Resolver
		DataSources struct {
			Dynamo []DynamoDataSource
		}
	}

	// AppSync represents the AppSync specific parts of the terraform
	AppSync struct {
		Name string
	}

	Resolver interface{}

	// DynamoDataSource declares a dynamoDB table to be used as
	// an appsync datasource
	DynamoDataSource struct {
		Name          string
		Identifier    string // The lowercased and stripped Name
		HashKey       *DynamoKey
		SortKey       *DynamoKey
		DisableBackup bool
	}

	// DynamoKey represents a key field in a dynamo data source. It defines
	// the name of the field as well as the dynamo data type.
	DynamoKey struct {
		Field string
		Type  string
	}
)

// NewDynamoKey parses a raw field definition (`<field>[:<type>]`) and returns
// a new poplulated DynamoKey
func NewDynamoKey(field string) *DynamoKey {
	if field == "" {
		return nil
	}
	dk := &DynamoKey{
		Field: field,
		Type:  "S",
	}
	if strings.Contains(field, ":") {
		fieldAndType := strings.Split(field, ":")
		dk.Field = fieldAndType[0]
		dk.Type = fieldAndType[1]
	}
	return dk
}

// New returns an empty terraform
func New() *Terraform {
	return &Terraform{
		AppSync:   AppSync{},
		Resolvers: []Resolver{},
	}
}

// NewFromManifest creates a terraform from a parsed manifest
func NewFromManifest(m *manifest.Manifest) (*Terraform, error) {

	tf := New()

	if m.APINameSuffix == "" {
		return nil, errors.New("manifest does not specify apiNameSuffix")
	}
	tf.AppSync.Name = m.APINameSuffix

	// Import data sources from the manifest
	for dsName, ds := range m.DataSources {
		source := ds.GetSource()
		switch t := source.(type) {
		case *manifest.DynamoSource:
			dds := DynamoDataSource{
				Name:          dsName,
				Identifier:    stripNameToIdentifier(dsName),
				DisableBackup: t.DisableBackup,
				HashKey:       NewDynamoKey(t.HashKey),
				SortKey:       NewDynamoKey(t.SortKey),
			}
			tf.DataSources.Dynamo = append(tf.DataSources.Dynamo, dds)
		case *manifest.RDSSource:
			log.Println("TERRAFORM NOT SUPPORTING RDS YET")
		default:
			return nil, fmt.Errorf("data source type not supported yet: %T", t)
		}

		// switch {
		// case ds.Dynamo != nil: // Dynamo
		// 	dds := DynamoDataSource{
		// 		Name:          dsName,
		// 		Identifier:    stripNameToIdentifier(dsName),
		// 		DisableBackup: ds.Dynamo.DisableBackup,
		// 		HashKey:       NewDynamoKey(ds.Dynamo.HashKey),
		// 		SortKey:       NewDynamoKey(ds.Dynamo.SortKey),
		// 	}
		// 	tf.DataSources.Dynamo = append(tf.DataSources.Dynamo, dds)
		// case ds.RDS

		// }
	}

	// TODO import parts
	// ...

	return tf, nil
}

func stripNameToIdentifier(name string) string {
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	return strings.ToLower(name)
}

// Marshal marshals the schema
func (tf *Terraform) Marshal() ([]byte, error) {
	// TODO Can we sensibly marshal as we really want this in chunks?
	return nil, nil
}

// Write outouts the generated terraform to a given io.Writer
func (tf *Terraform) Write(w io.Writer) error {
	return tmplTerraform.Execute(w, tf)
}
