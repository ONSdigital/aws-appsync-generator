package graphql

import (
	"fmt"
	"regexp"
)

var (
	reAWSScalars     = regexp.MustCompile(`(AWS(Date(Time)?|Time(stamp)?|Email|URL|Phone|IPAddress|JSON))`)
	reGraphQLScalars = regexp.MustCompile(`(ID|String|Int|Float|Boolean)`)
)

type (
	// Object represents a graphql object type definition
	Object struct {
		Name   string   `yaml:"name"`
		Fields []*Field `yaml:"fields"`
	}

	// FilterObject is generated from a base object and used to filter queries
	FilterObject Object

	// FilterObjectList is a list type containing FilterObjects
	FilterObjectList []*FilterObject

	// InputObject is generated from a base object and used for mutation queries
	InputObject Object

	// InputObjectList is a list type containing InputObjects
	InputObjectList []*InputObject
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

// NewInputFromObject creates a new input object type from an
// existing Object definition.
// Fields that can't be (currently) filtered (custom types and embedded objects)
// are omitted. Non-standard, but mappable, types are translated to simple
// scalar types
func NewInputFromObject(o *Object, action string) (*InputObject, error) {
	prefix := ""
	switch action {
	case ActionInsert:
		prefix = "Create"
	case ActionUpdate:
		prefix = "Update"
	default:
		return nil, fmt.Errorf("invalid action type for input object '%s': %s", o.Name, action)
	}
	io := &InputObject{
		Name: fmt.Sprintf("%s%sInput", prefix, o.Name),
	}
	for _, f := range o.Fields {
		fieldTypeName := f.Type.Name

		// Allow overriding of types to be different for the input object
		if f.InputType != nil {
			fieldTypeName = f.InputType.Name
		}

		// If this is an insert (create) object, then we omit any "id" field
		// with a type "ID!" as it should be autogenerated
		if action == ActionInsert && f.Name == "id" && f.Type.Name == "ID" && f.Type.NonNullable {
			continue
		}

		// TODO Deal with enum fields that should map as strings

		io.Fields = append(io.Fields, &Field{
			Name: f.Name,
			Type: &FieldType{
				Name:        fieldTypeName,
				IsList:      f.Type.IsList,
				NonNullable: false,
			},
		})
	}

	return io, nil
}

// AddInputFromObject adds a new filter object definition to the schema
func (s *Schema) AddInputFromObject(o *Object, action string) error {
	io, err := NewInputFromObject(o, action)
	if err != nil {
		return err
	}

	if s.InputObjects == nil {
		s.InputObjects = make(InputObjectList, 0, 1)
	}
	s.InputObjects = append(s.InputObjects, io)
	return nil
}
