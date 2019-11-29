package schema

type (
	// Query represents the definition for a graphql query
	Query struct {
		Name     string    `yaml:"name"`
		Type     string    `yaml:"type"`
		Resolver *Resolver `yaml:"resolver"`
	}
)
