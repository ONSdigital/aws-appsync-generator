package terraform

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"time"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
)

var (
	// OutputPrefix is prepended to the start of each generated file
	// name. This is used to help discriminate generated files from
	// manual files in an output location.
	OutputPrefix = "_"

	// OutputPath is the root path to write generated files to
	OutputPath = "./generated"

	// TemplatePath is where to find the terraform templates
	TemplatePath = "./templates"
)

var funcMap = template.FuncMap{
	"now": func() string {
		return time.Now().String()
	},
}

func Generate(m *manifest.Manifest) error {

	if err := cleanOutput(); err != nil {
		return errors.Wrap(err, "failed to clear output folder")
	}

	resolvers := []resolverData{}

	// Gather resolvers from user defined objects
	for _, o := range m.Objects {
		for _, f := range o.Fields {
			if r := f.Resolver; r != nil {
				resolvers = append(resolvers, resolverData{
					FieldName:        f.Name,
					ParentObjectName: o.Name,
					Action:           r.Action,
					KeyFields:        r.KeyFields,
					Returns:          r.Returns(),
					DataSourceType:   "dynamo", // TODO
					DataSourceName:   "",       // TODO
				})
			}
		}
	}

	for _, r := range resolvers {
		log.Printf("Write terraform: %s\n", r.FileName())
		if err := r.Write(); err != nil {
			return err
		}
	}
	return nil
}

// Generate creates all the terraform output to create the
// infrastructure for a manifest
func Generate2(m *manifest.Manifest) error {

	// Clean the output folder
	if err := cleanOutput(); err != nil {
		return errors.Wrap(err, "failed to clear output folder")
	}

	resolvers := []*resolverDefinition{}

	for _, o := range m.Objects {
		for _, f := range o.Fields {
			if f.Resolver != nil {
				resolvers = append(resolvers, &resolverDefinition{
					ParentName:  o.Name,
					FieldName:   f.Name,
					Resolver:    f.Resolver,
					ParentField: &f,
					// Action:     f.Resolver.Action,
				})
			}
		}
	}

	for _, f := range m.Queries {
		if f.Resolver == nil {
			return fmt.Errorf("query field (%s) must have a resolver", f.Name)
		}
		resolvers = append(resolvers, &resolverDefinition{
			ParentName:  "Query",
			FieldName:   f.Name,
			Resolver:    f.Resolver,
			ParentField: &f,
			// Action:     f.Resolver.Action,
		})
	}

	for _, f := range m.Mutations {
		if f.Resolver == nil {
			return fmt.Errorf("query field (%s) must have a resolver", f.Name)
		}
		resolvers = append(resolvers, &resolverDefinition{
			ParentName:  "Mutation",
			FieldName:   f.Name,
			Resolver:    f.Resolver,
			ParentField: &f,
			// Action:     f.Resolver.Action,
		})
	}

	for _, r := range resolvers {
		log.Printf("Write terraform: %s\n", r.FileName())
		if err := r.Write(); err != nil {
			return err
		}
	}

	return nil
}

func cleanOutput() error {
	d, err := os.Open(OutputPath)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(OutputPath, name))
		if err != nil {
			return err
		}
	}
	return nil
}
