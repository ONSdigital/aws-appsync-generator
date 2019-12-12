package graphql

// DataSource (string)
type DataSource string

// Constants for known datasource types
const (
	DataDynamo DataSource = "dynamo"
	DataSQL    DataSource = "sql"
)

type (
	// Source is a datasource configuration
	Source struct {
		Name string     `yaml:"name"`
		Type DataSource `yaml:"type"`
	}
)
