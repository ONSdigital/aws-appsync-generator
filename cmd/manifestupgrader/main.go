package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	app := &cli.App{
		Name:  "upgrade",
		Usage: "upgrade a v1 manifest to v2",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run generator",
				Action: runCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "manifest",
						Aliases:  []string{"m"},
						Usage:    "manifest file defining the api",
						EnvVars:  []string{"MANIFEST"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "name of the api",
						EnvVars:  []string{"NAME"},
						Required: true,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type (
	// A temporary representation of the v1 schema for unmarhsaling to.
	oldManifest struct {
		Enums []struct {
			Name   string
			Values []string
		}
		DataSources map[string]struct {
			Name   string
			Dynamo *struct {
				Backup  bool `yaml:"backup"`
				HashKey struct {
					Name string
					Type string
				} `yaml:"hash_key"`
				SortKey struct {
					Name string
					Type string
				} `yaml:"sort_key"`
			}
		} `yaml:"sources"`
		Objects []struct {
			Name   string
			Fields []struct {
				Name     string
				Type     interface{}
				Resolver *struct {
					Action    string
					Type      interface{}
					Source    string
					KeyFields []struct {
						Name   string
						Parent string
					} `yaml:"keyFields"`
				}
			}
		}
	}
)

func runCommand(c *cli.Context) error {

	if c.String("manifest") == "" {
		return errors.New("missing mandatory 'manifest' parameter")
	}

	if c.String("name") == "" {
		return errors.New("missing mandatory 'name' parameter")
	}

	log.Println("Converting to manifest version v2.0.0")

	// Assuming the file isn't too big
	oldData, err := ioutil.ReadFile(c.String("manifest"))
	if err != nil {
		return errors.Wrap(err, "failed to read old manifest")
	}

	// Marshal in the old schema
	var old oldManifest
	err = yaml.Unmarshal(oldData, &old)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal old manifest")
	}
	newManifest := manifest.Manifest{
		Version:       "v2.0.0",
		APINameSuffix: c.String("name"),
		Enums:         make(map[string]manifest.Enum),
		DataSources:   make(map[string]manifest.DataSource),
		Objects:       make(map[string]manifest.Object),
	}

	// Import enums
	for _, enum := range old.Enums {
		log.Println("Importing enum:", enum.Name)
		newManifest.Enums[enum.Name] = enum.Values
	}

	// Import data sources
	for _, ds := range old.DataSources {
		log.Println("Importing datasource:", ds.Name)
		if dds := ds.Dynamo; dds != nil {
			newManifest.DataSources[ds.Name] = manifest.DataSource{
				Dynamo: &manifest.DynamoSource{
					DisableBackup: !dds.Backup,
					HashKey:       formatDynamoKey(dds.HashKey.Name, dds.HashKey.Type),
					SortKey:       formatDynamoKey(dds.SortKey.Name, dds.SortKey.Type),
				},
			}
		}
	}

	// Import objects
	for _, obj := range old.Objects {
		log.Println("Importing object:", obj.Name)
		object := make(manifest.Object, 0)
		for _, field := range obj.Fields {
			log.Println("└──", field.Name)
			fieldName := field.Name
			// Type assert the type! (I was far too clever in the previous
			// version! - it was pretty cool but probably overkill for some
			// syntactic sugar)
			// This deals with the Type coming out as an interface{} or an
			// []interface{} rather than a concrete type.
			switch t := field.Type.(type) {
			case string:
				fieldName = fmt.Sprintf("%s:%s", fieldName, t)
			case []interface{}:
				switch tl := t[0].(type) {
				case string:
					fieldName = fmt.Sprintf("%s:[%s]", fieldName, tl)
				}
			default:
				// Nothing to do! Can leave the name just as it
				// is and not append any default type to it
			}

			// TODO resolvers
			// ...
			if field.Resolver != nil {
				log.Println("|   └── (+resolver)", field.Resolver.Action)
				// TODO ...
			}

			object = append(object, manifest.Field{Name: fieldName})
		}
		newManifest.Objects[obj.Name] = object
	}

	out, err := yaml.Marshal(newManifest)
	if err != nil {
		return err
	}

	// Output the converted manifest
	fmt.Println(string(out))
	return nil
}

func formatDynamoKey(name, dynamoType string) string {
	if dynamoType != "" {
		return name + ":" + dynamoType
	}
	return name
}
