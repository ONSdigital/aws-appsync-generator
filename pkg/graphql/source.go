package graphql

type DataSource string

const (
	DataDynamo DataSource = "dynamo"
	DataSQL    DataSource = "sql"
)

type (
	Source struct {
		Name string     `yaml:"name"`
		Type DataSource `yaml:"type"`
	}
)
