package manifest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v2"
)

var supportedMajorManifestVersion = "v2"
var reVersion = regexp.MustCompile("version: *(v.*)")

// TODO
// - Validate schema to ensure mandatory elements are not missing
//	- e.g. hashKey in dynamo data sources

type (

	// Manifest is the top level representation of a manifest configuraion that
	// has been parsed
	Manifest struct {
		// Meta block
		// Contains info such as the version of this schema in use and
		// the name of the api to create
		Version       string `yaml:"version"`
		APINameSuffix string `yaml:"apiNameSuffix"`

		// GraphQL schema type definitions
		Objects   map[string]Object   `yaml:"objects,omitempty"`
		Enums     map[string]Enum     `yaml:"enums,omitempty"`
		Queries   map[string]Resolver `yaml:"queries,omitempty"`
		Mutations map[string]Resolver `yaml:"mutations,omitempty"`

		// Data source definitions
		DataSources map[string]DataSource `yaml:"sources,omitempty"`

		parsed bool
	}

	// Object represents a graphql schema object
	Object []Field

	// Enum represents a graphql enum type
	Enum []string

	// Resolver represents a graphql resolver
	Resolver struct {
		Action string `yaml:"action"`
	}

	// Field represent a graphql schema object field type
	Field struct {
		Name             string `yaml:"field"`
		ExcludeFromInput bool   `yaml:"excludeFromInput,omitempty"`
	}
)

// IsNonNullable returns whether a given field has been marked as non-nullable
func (f *Field) IsNonNullable() bool {
	return strings.HasSuffix(f.Name, "!")
}

// Parse reads a given manifest and returns a new `Manifest`
func Parse(data []byte) (*Manifest, error) {

	// Check the version of the manifest before we go any further.
	// Gives us a quick(ish - as long as it's defined at all!) bailout
	// if we've been given a schema in the wrong format before throwing lots of
	// horrible parsing errors.
	if err := CheckVersion(data); err != nil {
		return nil, err
	}

	var m Manifest
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse manifest")
	}

	// Ensure we've been given a name
	if m.APINameSuffix == "" {
		return nil, errors.New("must supply apiNameSuffix")
	}

	err = m.ValidateDataSources()
	var ve *ValidationError
	if errors.As(err, &ve) {

		return nil, errors.Unwrap(err)
	}

	m.parsed = true
	return &m, nil
}

// ValidationError wraps errors returned from validating a manifest
type ValidationError struct {
	Errors []error
}

func (ve ValidationError) Error() string {
	return "errors found whilst validating"
}

// ValidateDataSources runs through all declared datasources in the manifest and
// attempts to check their configurations. Makes no attempt to actually see if
// datasources are real databases etc!
func (m *Manifest) ValidateDataSources() error {
	ve := ValidationError{
		Errors: []error{},
	}
	for dsName, ds := range m.DataSources {
		if err := ds.Validate(); err != nil {
			ve.Errors = append(ve.Errors, errors.Wrapf(err, "validation error on datasource: %s", dsName))
		}
	}
	if len(ve.Errors) > 0 {
		return ve
	}
	return nil
}

// CheckVersion verifies that given manifest data is of a compatible version
func CheckVersion(data []byte) error {
	r := bufio.NewReader(bytes.NewReader(data))
	for {
		line, err := r.ReadBytes('\n')
		if err == io.EOF {
			return errors.New("no valid 'version' found in manifest - must be in format 'v1.0.0'")
		}
		if err != nil {
			return err
		}

		matches := reVersion.FindSubmatch(line)
		if len(matches) > 0 {
			version := string(matches[1])
			if semver.Major(version) != supportedMajorManifestVersion {
				return fmt.Errorf("manifest must be version >= 2.0.0 and < 3.0.0, got %s", version)
			}
			break
		}
	}
	return nil
}

// GetAttributeName takes an attribute in the form `<name[:type]>` and returns
// the `name` portion
func GetAttributeName(attr string) string {
	if strings.Contains(attr, ":") {
		return strings.Split(attr, ":")[0]
	}
	return attr
}

// GetAttributeType takes an attribute in the form `<name[:type]` and returns
// the `type` portion. If none was specified, then the default given by `def`
// is returned instead
func GetAttributeType(attr, def string) string {
	if strings.Contains(attr, ":") {
		return strings.Split(attr, ":")[1]
	}
	return def
}
