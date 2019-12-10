package graphql

type (
	// Object represents a graphql object type definition
	Object struct {
		Name   string   `yaml:"name"`
		Fields []*Field `yaml:"fields"`
	}
)
