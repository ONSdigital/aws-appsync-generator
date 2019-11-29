package schema

type (
	// Mutation is the definition for a mutation query
	Mutation struct {
		Query `yaml:",inline"` // Using the inline allows the unmarshaller to correctly reach the inner fields

		// KeyFields are the identifiers to be used to pull back the mutated record
		KeyFields []*Field `yaml:"keyFields"`
	}
)
