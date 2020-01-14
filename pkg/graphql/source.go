package graphql

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type (
	// Source is a datasource configuration
	Source struct {
		Name   string        `yaml:"name"`
		Dynamo *DynamoSource `yaml:"dynamo"`
		SQL    *SQLSource    `yaml:"sql"`

		// Set automatically
		Type string
	}

	// DynamoSource represents a dynamo db data source
	DynamoSource struct {
		HashKey string `yaml:"hash_key"`
		SortKey string `yaml:"sort_key,omitempty"`
	}

	// SQLSource represents a sql based db data source
	SQLSource struct {
		PrimaryKey string `yaml:"primary_key"`
		// TODO other fields
	}

	unmarshalSource Source
)

var reSupportedDataSourceTypes = regexp.MustCompile(`(dynamo|aurora)`)

// UnmarshalYAML satisfies the custom unmarshaler interface for go-yaml. It is
// called automatically by the YAML unmarshal.
func (ds *Source) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var u unmarshalSource
	if err := unmarshal(&u); err != nil {
		return err
	}

	*ds = Source(u)

	if ds.Name == "" {
		return errors.New("datasource does not declare a name")
	}

	switch {
	case ds.Dynamo != nil:
		ds.Type = "dynamo"
	case ds.SQL != nil:
		ds.Type = "sql"
	default:
		return errors.New("must specify a support data source type")
	}
	return nil
}

// GenerateBytes renders the datasource ready to be written to the output stream
func (ds *Source) GenerateBytes() ([]byte, error) {
	generated := bytes.Buffer{}

	t, err := template.New(ds.Name).Funcs(funcMap).Parse(sourceTemplate)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(&generated, ds); err != nil {
		return nil, err
	}
	return generated.Bytes(), nil
}

// OutputName returns the file name to be written for the data source
func (ds *Source) OutputName() string {
	return strings.ToLower(fmt.Sprintf("_datasource_%s_%s.tf", ds.Type, ds.Name))
}
