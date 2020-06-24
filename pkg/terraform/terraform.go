package terraform

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/ONSdigital/aws-appsync-generator/pkg/mapping"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

type (
	// Terraform represents the overall data to populate the terraform
	Terraform struct {
		AppSync           AppSync
		Resolvers         []*Resolver
		ResolverTemplates mapping.Templates
		// DataSources       struct {
		// 	Dynamo []DynamoDataSource
		// }
		DataSources map[string]map[string]DataSource

		TemplatePath string
	}

	// AppSync represents the AppSync specific parts of the terraform
	AppSync struct {
		Name string
	}

	// DataSource is the interface for all data source types
	DataSource interface{}

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
		Resolvers: []*Resolver{},
		DataSources: map[string]map[string]DataSource{
			"dynamo": make(map[string]DataSource),
			"sql":    make(map[string]DataSource),
			"lambda": make(map[string]DataSource),
		},
	}
}

// NewFromManifest creates a terraform from a parsed manifest
func NewFromManifest(m *manifest.Manifest, templates mapping.Templates) (*Terraform, error) {

	tf := New()
	tf.ResolverTemplates = templates

	// .....
	spew.Dump(m.Queries) // FIXME REMOVE
	// .....

	if err := tf.ImportFromManifest(m); err != nil {
		return nil, err
	}
	return tf, nil
}

// ImportFromManifest populates the terraform struct from the information in the
// given manifest
func (tf *Terraform) ImportFromManifest(m *manifest.Manifest) error {
	if m.APINameSuffix == "" {
		return errors.New("manifest does not specify apiNameSuffix")
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
			tf.DataSources["dynamo"][dsName] = dds
		case *manifest.RDSSource:
			log.Println("TERRAFORM NOT SUPPORTING RDS YET")
		default:
			return fmt.Errorf("data source type not supported yet: %T", t)
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

	resolvers, err := m.Resolvers()
	if err != nil {
		return errors.New("failed to retrieve resolvers from manifest")
	}
	for name, r := range resolvers {

		resolver := &Resolver{}
		if err := resolver.BuildDefinition(r); err != nil {
			return errors.Wrap(err, "error parsing resolver definition")
		}

		// 1. Get the associated data source for the resolver
		// 2. Find the type of the data source
		// 3. Use that to find the relevant sub-loaded template

		ds := tf.getDataSourceByName(r.Source)
		if ds == nil {
			return fmt.Errorf("no datasource with name '%s' decalred", r.Source)
		}

		dataSourceType := ""

		switch tp := ds.(type) {
		case DynamoDataSource:
			dataSourceType = "dynamo"
		default:
			return fmt.Errorf("unsupported datasource type: %s", tp)
		}

		t, ok := tf.ResolverTemplates[dataSourceType][r.Template]
		if !ok {

			for k, v := range tf.ResolverTemplates {
				log.Println(k)
				for k := range v {
					log.Println("  ", k)
				}
			}

			return fmt.Errorf("template %s/%s not loaded", dataSourceType, r.Template)
		}

		// TODO
		// Build the vars that need to be passed to the templates
		// Then actually pass to the templates

		var b bytes.Buffer

		if err = t.ExecuteTemplate(&b, "signature", nil); err != nil {
			return errors.Wrapf(err, "failed to create signature for resolver '%s'", name)
		}
		resolver.Signature = b.String()
		b.Reset()

		if err = t.ExecuteTemplate(&b, "request", nil); err != nil {
			return errors.Wrapf(err, "failed to create request for resolver '%s'", name)
		}
		resolver.Request = b.String()
		b.Reset()

		if err = t.ExecuteTemplate(&b, "request", nil); err != nil {
			return errors.Wrapf(err, "failed to create response for resolver '%s'", name)
		}
		resolver.Response = b.String()
		b.Reset()

		tf.Resolvers = append(tf.Resolvers, resolver)
	}
	return nil
}

func (tf *Terraform) getDataSourceByName(name string) DataSource {
	for _, sources := range tf.DataSources {
		if source, ok := sources[name]; ok {
			return source
		}
	}
	return nil
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
