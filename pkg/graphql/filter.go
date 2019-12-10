package graphql

import (
	"fmt"
	"regexp"
)

type (
	// FilterObject is generated from a base object and used to filter queries
	FilterObject Object

	// FilterObjectList is a list of FilterObject
	FilterObjectList []*FilterObject
)

var (
	reAWSScalars     = regexp.MustCompile(`(AWS(Date(Time)?|Time(stamp)?|Email|URL|Phone|IPAddress|JSON))`)
	reGraphQLScalars = regexp.MustCompile(`(ID|String|Int|Float|Boolean)`)
)

// NewFilterFromObject creates a new filter object input type from an
// existing Object definition.
// Fields that can't be (currently) filtered (custom types and embedded objects)
// are omitted. Non-standard, but mappable, types are translated to simple
// scalar types
func NewFilterFromObject(o *Object) *FilterObject {
	fo := &FilterObject{
		Name: o.Name + "Filter",
	}
	for _, f := range o.Fields {
		fieldTypeName := ""
		switch {
		case reAWSScalars.MatchString(f.Type.Name):
			fieldTypeName = "TableStringFilterInput"
		case reGraphQLScalars.MatchString(f.Type.Name):
			fieldTypeName = fmt.Sprintf("Table%sFilterInput", f.Type.Name)
		default:
			// Skip field
			continue
		}

		// TODO Deal with enum fields that should map as strings

		fo.Fields = append(fo.Fields, &Field{
			Name: f.Name,
			Type: &FieldType{
				Name:        fieldTypeName,
				IsList:      false,
				NonNullable: false,
			},
		})
	}
	return fo
}

// AddFilterFromObject adds a new filter object definition to the schema
func (s *Schema) AddFilterFromObject(o *Object) {
	if s.FilterObjects == nil {
		s.FilterObjects = make(FilterObjectList, 0, 1)
	}
	s.FilterObjects = append(s.FilterObjects, NewFilterFromObject(o))
}
