package graphql

type (
	// Query represents the definition for a graphql query
	//
	//   Query {
	//   - name: <query_name>
	//     type: <return_type>
	//    resolver:
	//       action: <get|list|insert|update>
	//       source: <source_name> # optional
	//   ...
	//   }
	//
	Query struct {
		Name     string     `yaml:"name"`
		Type     *FieldType `yaml:"type"`
		Resolver *Resolver  `yaml:"resolver"`
	}
)
