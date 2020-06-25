package terraform

// import (
// 	"fmt"

// 	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
// )

// type (
// 	// Resolver is the definition of a resolver that is to be output to
// 	// a terraform script
// 	Resolver struct {
// 		Identifier string
// 		Signature  string //
// 		Source     string
// 		ParentType string //
// 		Request    string //
// 		Response   string //
// 		ReturnType string //
// 		FieldName  string //

// 		DataSourceType       string
// 		DataSourceIdentifier string

// 		KeyFields []manifest.Field
// 	}
// )

// func (r *Resolver) String() string {
// 	return fmt.Sprintf("[%s] %s", r.ParentType, r.FieldName)
// }

// // ArgsSource returns "args" or "source" depending on whether the parent is a
// // custom object or query/mutation. Expects resolver to have already be populated
// // with a valid `ParentType`. If no type is set it will return `source` which
// // may not be what you want.
// func (r *Resolver) ArgsSource() string {
// 	switch r.ParentType {
// 	case "Query":
// 		return "args"
// 	case "Mutation":
// 		return "args"
// 	default:
// 		return "source"
// 	}
// }

// // KeyFieldJSONMap converts the `KeyFields` into a JSON formatted map suitable
// // for use in a VTL template
// func (r *Resolver) KeyFieldJSONMap() string {
// 	fl := make([]string, len(r.KeyFields))
// 	for i, f := range r.KeyFields {
// 		// Remove any "mandatory" mark as it doesn't make sense in the usage
// 		// contexts of the resultant map
// 		returnType := manifest.GetAttributeType(f.Name, "String")
// 		returnType = strings.TrimSuffix(returnType, "!")
// 		fl[i] = fmt.Sprintf(`"%s":"%s"`, manifest.GetAttributeName(f.Name), returnType)
// 	}
// 	return "{" + strings.Join(fl, ",") + "}"
// }

// // KeyFieldJSONList converts the `KeyFields` into a JSON formatted list of names
// // suitable for use in a VTL template
// func (r *Resolver) KeyFieldJSONList() string {
// 	fl := make([]string, len(r.KeyFields))
// 	for i, f := range r.KeyFields {
// 		fl[i] = fmt.Sprintf(`"%s"`, manifest.GetAttributeName(f.Name))
// 	}
// 	return "[" + strings.Join(fl, ",") + "]"
// }
