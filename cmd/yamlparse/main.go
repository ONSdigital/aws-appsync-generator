package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/ONSdigital/aws-appsync-generator/pkg/schema"
	"github.com/pkg/errors"

	flag "github.com/spf13/pflag"
)

// Path constants
const (
	GeneratedFilesPath = "generated"
	PublicSchemaName   = GeneratedFilesPath + "/schema.public.graphql"

	TemplateRoot = "templates"
)

var (
	manifest     = ""
	targetDBType = ""
)

func init() {
	flag.StringVarP(&manifest, "manifest", "m", "manifest.yml", "manifest file to parse")
	flag.StringVarP(&targetDBType, "target", "t", "", "target db - sql or dynamodb")
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

	s, err := schema.New(body)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse definition"))
	}

	// // Build the input objects - TODO refactor
	for _, in := range s.Inputs {

		// TODO - Assuming we have a base object for now. Not supporting
		// 		  arbitrary params for now

		// Populate the input object's fields. Can only do this once the
		// schema definition has been parsed
		if err := in.Populate(s); err != nil {
			log.Fatalf("unable to populate: %v", err) // TODO errors
		}
	}

	// GRAPHQL SCHEMA BUILD
	var generatedSchema *bytes.Buffer
	if generatedSchema, err = generatePublicSchema(s); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(generatedSchema.String())
	if err := write(PublicSchemaName, generatedSchema); err != nil {
		log.Fatalf("failed to output: %v", err)
	}

	for _, object := range s.Objects {
		for _, field := range object.Fields {

			// Object field resolvers
			if r := field.Resolver; r != nil {

				generated, resolverIdentifier, err := generateResolverTerraform(r, field.Name, object.Name, "source")
				if err != nil {
					log.Fatal(err) // TODO
				}

				if err := write(GeneratedFilesPath+"/terraform/"+resolverIdentifier+".tf", generated); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	log.Println("Generating query terraform")
	for _, query := range s.Queries {

		r := query.Resolver
		if r == nil {
			log.Fatal("Missing resolver for query", query.Name)
		}

		generated, resolverIdentifier, err := generateResolverTerraform(r, query.Name, "Query", "args")
		if err != nil {
			log.Fatal(err) // TODO
		}

		if err := ioutil.WriteFile(GeneratedFilesPath+"/terraform/"+resolverIdentifier+".tf", generated.Bytes(), 0644); err != nil {
			log.Fatal(err) // TODO
		}
	}

	log.Println("Generating mutation terraform")
	for _, mutation := range s.Mutations {

		// TODO Refactor. Need a more elegant way of getting the input fields
		// without necessarily needing to explicitly promote them here. The
		// resolver should likely do this automatically
		// (need to look at replacing Resolver with interface)

		r := mutation.Resolver
		if r == nil {
			log.Fatal("Missing resolver for mutation", mutation.Name)
		}

		// Get the input object for the mutation and attach its params to
		// the resolver
		in, err := s.Input(r.Input)
		if err != nil {
			log.Fatal("bad base input:", err) // TODO errors
		}

		// Promote the params
		r.Params = in.Params

		generated, resolverIdentifier, err := generateResolverTerraform(r, mutation.Name, "Mutation", "args")
		if err != nil {
			log.Fatal(err) // TODO
		}

		if err := ioutil.WriteFile(GeneratedFilesPath+"/terraform/"+resolverIdentifier+".tf", generated.Bytes(), 0644); err != nil {
			log.Fatal(err) // TODO
		}
	}

	log.Println("[Generation completed]")
}

var templateCache map[string]*template.Template

// TemplateData is passed to a new template to populate it
type TemplateData struct {
	ResolverIdentifier string
	Parent             string
	FieldName          string
	Resolver           *schema.Resolver
	ArgSource          string
}

func generateResolverTerraform(r *schema.Resolver, name, parent, source string) (*bytes.Buffer, string, error) {

	resolverIdentifier := strings.ToLower("_" + parent + "-" + name)
	generated := bytes.Buffer{}

	requestTemplate := r.Action
	if r.CustomRequestTemplateName != "" {
		requestTemplate = r.CustomRequestTemplateName
	}

	responseTemplate := r.Action
	if r.CustomResponseTemplateName != "" {
		responseTemplate = r.CustomResponseTemplateName
	}

	templateIdentifier := fmt.Sprintf("%s-%s", requestTemplate, responseTemplate)

	log.Printf("(resolver) [%s] %s\n", parent, resolverIdentifier)

	if templateCache == nil {
		templateCache = make(map[string]*template.Template)
	}

	t, ok := templateCache[templateIdentifier]
	if !ok {
		log.Printf("Template '%s' does not exist in cache - creating", templateIdentifier)
		// Template combination does not already exist
		files := []string{
			TemplateRoot + "/terraform/resolver.tmpl",
			fmt.Sprintf("%s/resolvers/%s/request/%s.tmpl", TemplateRoot, targetDBType, requestTemplate),
			fmt.Sprintf("%s/resolvers/%s/response/%s.tmpl", TemplateRoot, targetDBType, responseTemplate),
		}
		var err error
		t, err = template.New("resolver.tmpl").Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return nil, resolverIdentifier, errors.Wrapf(err, "failed to generate resolver '%s'", resolverIdentifier)
		}

		templateCache[templateIdentifier] = t
	}

	if err := t.Execute(&generated, TemplateData{resolverIdentifier, parent, name, r, source}); err != nil {
		return nil, resolverIdentifier, errors.Wrapf(err, "error executing template for resolver '%s'", resolverIdentifier)
	}

	return &generated, resolverIdentifier, nil
}

func write(path string, content *bytes.Buffer) error {
	return ioutil.WriteFile(path, content.Bytes(), 0644)
}

var funcMap = template.FuncMap{
	"now": func() string {
		return time.Now().String()
	},
	"joinFieldNames": func(s []*schema.Field) string {
		names := make([]string, 0, len(s))
		for _, f := range s {
			names = append(names, f.Name)
		}
		return strings.Join(names, ",")
	},
}

func generatePublicSchema(s *schema.Schema) (*bytes.Buffer, error) {
	generated := bytes.Buffer{}

	templateName := "schema.public.graphql.tmpl"

	t, _ := template.New(templateName).Funcs(funcMap).ParseFiles("templates/" + templateName)
	if err := t.Execute(&generated, s); err != nil {
		return nil, err
	}
	return &generated, nil
}
