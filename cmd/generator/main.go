package main

import (
	"io/ioutil"
	"log"

	"github.com/ONSdigital/aws-appsync-generator/pkg/manifest"
	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

var (
	manifestFile = ""
)

func init() {
	flag.StringVarP(&manifestFile, "manifest", "m", "manifest.yml", "manifest file to parse")
	// flag.StringVarP(&graphql.GeneratedFilesPath, "output", "o", graphql.GeneratedFilesPath, "path to output generated files to (CAUTION: will be emptied before write!)")
	flag.Parse()
}

func main() {
	body, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read manifest '%s'", manifestFile))
	}

	m, err := manifest.New(body)
	if err != nil {
		log.Fatal(err)
	}

	err = schema.Generate(m)
	if err != nil {
		log.Fatal(err)
	}

	// err = terraform.Generate(m)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// spew.Dump(m)

}
