package schema

import "fmt"

// Object represents an object type in a graphql schema
type Object struct {
	// The name of the object
	Name string

	// If set true, an associated Input object will also be created
	// in the public schema
	// MakeInput bool `yaml:"makeInput"`

	// Sub
	Fields    []*Field
	Mutations interface{} // TODO

	// Generated
	fieldMap map[string]*Field
}

// Field returns the information for a particular field in an object
// definition
func (o *Object) Field(name string) (*Field, error) {
	if f, ok := o.fieldMap[name]; ok {
		return f, nil
	}
	return nil, fmt.Errorf("field '%s' not found", name)
}
