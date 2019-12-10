package graphql

type (
	// Enum represents a graphql enumeration definition
	Enum struct {
		Name   string   `yaml:"name"`
		Values []string `yaml:"values"`
	}
)
