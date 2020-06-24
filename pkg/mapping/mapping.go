// Package mapping contains functionality for discovering and loading custom
// resolver mapping templates
package mapping

import "text/template"

type (
	templateMap map[string]*template.Template

	// Templates contains a lookup mapping of
	Templates map[string]templateMap
)

// New returns an empty initialised templates map
func New() Templates {
	return Templates{
		"dynamo": make(templateMap),
		"sql":    make(templateMap),
		"lambda": make(templateMap),
		// TODO other source types
	}
}
