// Package mapping contains functionality for discovering and loading standard
// and custom resolver mapping templates
package mapping

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type (
	Templates map[string]map[string]*Template

	// Templates contains a lookup mapping of resolver templates
	// Templates struct {
	// 	templates map[string]map[string]*Template
	// }

	// Template represents the parsed and compiled stringified data for
	// a particular template
	Template struct {
		Signature string
		Request   string
		Response  string
		// Signature *template.Template
		// Request   *template.Template
		// Response  *template.Template
		Template *template.Template
	}
)

// New returns an empty initialised templates map
func New() Templates {
	// return &Templates{
	// 	templates: make(map[string]map[string]*Template),
	// }
	return make(map[string]map[string]*Template)
}

// Import will import a set of mapping templates from the given folder
func (tm Templates) Import(folderName string) error {

	// We're expecting subfolders named after datasource types
	// at this level - dynamo, sql etc
	// Skip anything that isn't a directory type
	folders, err := ioutil.ReadDir(folderName)
	if err != nil {
		return errors.Wrap(err, "failed to read template folder")
	}

	for _, subfolder := range folders {
		if !subfolder.IsDir() {
			continue
		}
		dataSourceType := strings.ToLower(subfolder.Name())

		files, err := ioutil.ReadDir(filepath.Join(folderName, dataSourceType))
		if err != nil {
			return errors.Wrap(err, "failed to read template datasource folder")
		}

		for _, file := range files {

			if !strings.HasSuffix(file.Name(), ".tmpl") {
				continue
			}

			log.Printf("Importing: %s/ (%s) %s", folderName, dataSourceType, file.Name())

			templateName := strings.TrimSuffix(file.Name(), ".tmpl")
			t, err := template.New(templateName).ParseFiles(
				filepath.Join(
					folderName,
					dataSourceType,
					file.Name(),
				),
			)
			if err != nil {
				return errors.Wrap(err, "failed to parse template")
			}

			tm.Add(dataSourceType, templateName, t)
		}
	}
	return nil
}

// Get fetches a template from the mapping lookup
func (tm Templates) Get(dataSourceType, templateName string) (*Template, error) {
	var ok bool
	var t *Template

	if _, ok = tm[dataSourceType]; !ok {
		return nil, fmt.Errorf("no mappings found under '%s'", dataSourceType)
	}

	if t, ok = tm[dataSourceType][templateName]; !ok {
		return nil, fmt.Errorf("no mapping '%s' found under '%s'", templateName, dataSourceType)
	}
	return t, nil
}

// Add executes a given template and adds it into the mapping at the given point
func (tm Templates) Add(dataSourceType, templateName string, t *template.Template) error {
	tmpl := &Template{
		Template: t,
	}
	// var b bytes.Buffer
	// var err error

	// if err = t.ExecuteTemplate(&b, "signature", nil); err != nil {
	// 	return errors.Wrapf(err, "failed to create signature for resolver '%s'", templateName)
	// }
	// tmpl.Signature = b.String()
	// b.Reset()

	// if err = t.ExecuteTemplate(&b, "request", nil); err != nil {
	// 	return errors.Wrapf(err, "failed to create request for resolver '%s'", templateName)
	// }
	// tmpl.Request = b.String()
	// b.Reset()

	// if err = t.ExecuteTemplate(&b, "request", nil); err != nil {
	// 	return errors.Wrapf(err, "failed to create response for resolver '%s'", templateName)
	// }
	// tmpl.Response = b.String()
	// b.Reset()

	if _, ok := tm[dataSourceType]; !ok {
		tm[dataSourceType] = make(map[string]*Template)
	}
	tm[dataSourceType][templateName] = tmpl

	return nil
}
