package graphql

import (
	"errors"
	"fmt"
	"regexp"
)

type (
	// Source is a datasource configuration
	Source struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
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

	if !reSupportedDataSourceTypes.MatchString(ds.Type) {
		return fmt.Errorf("datasource '%s' has unsupported type: %s", ds.Name, ds.Type)
	}
	return nil
}
