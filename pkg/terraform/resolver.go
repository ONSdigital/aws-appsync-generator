package terraform

import "github.com/ONSdigital/aws-appsync-generator/pkg/manifest"

type (
	// Resolver is the definition of a resolver that is to be output to
	// a terraform script
	Resolver struct {
		Signature string
		Source    string

		Request  string
		Response string
	}
)

// BuildDefinition takes a description of a resolver as given from a manifest
// and builds this terraform resolver definition from it
func (r *Resolver) BuildDefinition(mr *manifest.Resolver) error {

	// TODO ...

	return nil
}
