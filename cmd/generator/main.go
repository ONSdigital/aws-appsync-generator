package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ONSdigital/aws-appsync-generator/pkg/graphql"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

var (
	manifest     = ""
	targetDBType = ""
)

func init() {
	flag.StringVarP(&manifest, "manifest", "m", "manifest.yml", "manifest file to parse")
	flag.StringVarP(&targetDBType, "target", "t", "", "target db - sql or dynamo")
	flag.Parse()

	if targetDBType != "sql" && targetDBType != "dynamo" {
		fmt.Println("Target must be supplied and be one of 'sql' or 'dynamo'")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	body, err := ioutil.ReadFile(manifest)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to read manifest '%s'", manifest))
	}

	s, err := graphql.NewSchemaFromManifest(body, targetDBType)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse definition"))
	}

	if err := s.WriteAll(); err != nil {
		fmt.Println(err)
		for _, e := range s.Errors {
			fmt.Printf("(error) %v\n", e.Error())
		}
		fmt.Println("DONE (with errors)")
		os.Exit(1)
	}

	fmt.Println("DONE")
}
