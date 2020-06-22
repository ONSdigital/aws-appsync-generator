package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/ONSdigital/aws-appsync-generator/pkg/serverless"
	"github.com/ONSdigital/aws-appsync-generator/pkg/terraform"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	// Version is set by build flags
	Version = "0.0.0"
)

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
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "run generator",
				Action: runCommand,
			},
		},
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

	// spew.Dump(man)

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

	tf, err := terraform.NewFromManifest(man)
	if err != nil {
		return err
	}

	terraformName := "main.tf"
	if p := c.String("outpath"); p != "" {
		terraformName = path.Join(p, terraformName)
	}
	fTerraform, err := os.Create(terraformName)
	if err != nil {
		return err
	}
	tf.Write(fTerraform)
	fTerraform.Close()

	// serverless
	sless, err := serverless.NewFromManifest(man)
	if err != nil {
		return err
	}
	serverlessName := "serverless.yml"
	if p := c.String("outpath"); p != "" {
		serverlessName = path.Join(p, serverlessName)
	}
	fServerless, err := os.Create(serverlessName)
	if err != nil {
		return err
	}
	sless.Write(fServerless)
	fServerless.Close()

	return nil
}
