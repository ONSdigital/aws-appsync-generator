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

const defaultMappingTemplatesFolder = "mapping-templates"

// TODO
// - Validate schema to ensure mandatory elements are not missing
//	- e.g. hashKey in dynamo data sources

type (
	// ResolverMap is synatic sugar for the nested resolver map type
	ResolverMap map[string]map[string]*Resolver
)

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
		Objects   map[string]Object `yaml:"objects,omitempty"`
		Enums     map[string]Enum   `yaml:"enums,omitempty"`
		Queries   Object            `yaml:"queries,omitempty"`
		Mutations Object            `yaml:"mutations,omitempty"`

		// Data source definitions
		DataSources map[string]DataSource `yaml:"sources,omitempty"`

		path      string
		resolvers map[string]*Resolver
		parsed    bool

		// Resolvers mapped by parent (query, mutation or
		// object name) to field name
		resolversByParent ResolverMap
	}

	// Object represents a graphql schema object
	Object []Field

	// Enum represents a graphql enum type
	Enum []string

	// Field represent a graphql schema object field type
	Field struct {
		Name             string    `yaml:"field"`
		ExcludeFromInput bool      `yaml:"excludeFromInput,omitempty"`
		Resolver         *Resolver `yaml:"resolver,omitempty"`
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

	err = m.DiscoverResolvers()
	if err != nil {
		return nil, err
	}

	m.parsed = true
	return &m, nil
}

// Resolvers returns all resolvers currently discovered from the manifest. Can
// only be called once the manifest has been successfully parsed. Attempted to do
// so before will return an error
func (m *Manifest) Resolvers() (ResolverMap, error) {
	if !m.parsed {
		return nil, errors.New("manifest not parsed")
	}
	return m.resolversByParent, nil
}

// DiscoverResolvers finds all resolvers defined on fields, queries and mutations
// in the manifest and populates the `Resolvers` lookup.
// Resolvers are associated to their parent type. This will either be a "special"
// type (query / mutation) or the name of the parent object.
func (m *Manifest) DiscoverResolvers() error {
	m.resolvers = make(map[string]*Resolver)
	m.resolversByParent = make(map[string]map[string]*Resolver)

	m.resolversByParent["query"] = make(map[string]*Resolver)
	for _, q := range m.Queries {
		if q.Resolver != nil {
			q.Resolver.ParentType = "Query"
			q.Resolver.FieldName = q.Name
			m.resolversByParent["query"][q.Name] = q.Resolver
		}
	}

	m.resolversByParent["mutation"] = make(map[string]*Resolver)
	for _, mt := range m.Mutations {
		if mt.Resolver != nil {
			mt.Resolver.ParentType = "Mutation"
			mt.Resolver.FieldName = mt.Name
			m.resolversByParent["mutation"][mt.Name] = mt.Resolver
		}
	}

	for on, fields := range m.Objects {
		m.resolversByParent[on] = make(map[string]*Resolver)
		for _, f := range fields {
			if f.Resolver != nil {
				f.Resolver.ParentType = on
				f.Resolver.FieldName = f.Name
				m.resolversByParent[on][f.Name] = f.Resolver
			}
		}
	}

	return nil
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

// GetAttributeTypeStripped does the same as GetAttributeType, with the addition
// of stripped any mandatory flag (!) or list delimiters ([])
// e.g. [Animal!] -> Animal
func GetAttributeTypeStripped(attr, def string) string {
	if strings.Contains(attr, ":") {
		s := strings.Split(attr, ":")[1]
		s = strings.TrimPrefix(s, "[")
		s = strings.TrimSuffix(s, "]")
		s = strings.TrimSuffix(s, "!")
		return s
	}
	return def
}
