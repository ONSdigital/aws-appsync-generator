package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/ONSdigital/aws-appsync-generator/pkg/mapping"
	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/ONSdigital/aws-appsync-generator/pkg/terraform"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// Version is set by build flags
var Version = "0.0.0"

var dataSourceTypes = []string{"dynamo", "sql", "lambda"}

func main() {
	app := &cli.App{
		Name:    "appsyncgen",
		Version: Version,
		Usage:   "cli interface for running the Appsync generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "manifest",
				Aliases:  []string{"m"},
				Usage:    "manifest file defining the api",
				EnvVars:  []string{"MANIFEST"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "outpath",
				Aliases: []string{"o"},
				Usage:   "specify an output path",
			},
			&cli.StringFlag{
				Name:    "templates",
				Aliases: []string{"t"},
				Usage:   "specify a custom template folder",
			},
		},
		Action: runCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(c *cli.Context) error {

	manifestFile := c.String("manifest")

	body, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read manifest '%s'", manifestFile))
	}

	man, err := manifest.Parse(body)
	if err != nil {
		return err
	}

	s, err := schema.NewFromManifest(man)
	if err != nil {
		return err
	}

	schemaName := "schema.public.graphql"
	if p := c.String("outpath"); p != "" {
		schemaName = path.Join(p, schemaName)
	}
	fSchema, err := os.Create(schemaName)
	if err != nil {
		return err
	}
	s.Write(fSchema)
	fSchema.Close()

	templateLookup := mapping.New()

	// Standard mapping templates
	var templatePath string
	templatePath = "templates"
	if err := importTemplates(templateLookup, templatePath); err != nil {
		return errors.Wrap(err, "failed to import standard templates")
	}

	// Custom mapping templates
	templatePath = filepath.Join(filepath.Dir(manifestFile), "mapping-templates")
	if err := importTemplates(templateLookup, templatePath); err != nil {
		return errors.Wrap(err, "failed to import custom templates")
	}

	// terraform ------------
	tf, err := terraform.NewFromManifest(man, templateLookup)
	if err != nil {
		return err
	}

	log.Println("Generating terraform")
	terraformName := "main.tf"
	if p := c.String("outpath"); p != "" {
		terraformName = path.Join(p, terraformName)
	}
	log.Println("Writing terraform output to:", terraformName)
	fTerraform, err := os.Create(terraformName)
	if err != nil {
		return err
	}
	tf.Write(fTerraform)
	fTerraform.Close()

	// // serverless -----------
	// log.Println("Generating serverless")
	// sless, err := serverless.NewFromManifest(man, templates)
	// if err != nil {
	// 	return err
	// }
	// serverlessName := "serverless.yml"
	// if p := c.String("outpath"); p != "" {
	// 	serverlessName = path.Join(p, serverlessName)
	// }
	// fServerless, err := os.Create(serverlessName)
	// if err != nil {
	// 	return err
	// }
	// sless.Write(fServerless)
	// fServerless.Close()

	return nil
}

func importTemplates(lookup mapping.Templates, templatePath string) error {
	for _, dst := range dataSourceTypes {

		templateDataSourcePath := filepath.Join(templatePath, dst)

		if _, err := os.Stat(templateDataSourcePath); os.IsNotExist(err) {
			// Ignore missing datasource configurations
			continue
		}
		files, err := ioutil.ReadDir(templateDataSourcePath)
		if err != nil {
			return err
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".tmpl") {
				name := strings.TrimSuffix(file.Name(), ".tmpl")
				t, err := template.New(name).ParseFiles(
					filepath.Join(templateDataSourcePath, file.Name()),
				)
				if err != nil {
					return errors.Wrap(err, "failed to parse template")
				}
				lookup[dst][name] = t
			}
		}
	}
	return nil
}
