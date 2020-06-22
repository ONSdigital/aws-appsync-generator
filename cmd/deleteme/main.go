package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

func main() {
	y := []byte(`---
objects:
    Animal:
    - name: id:ID!
      excludeFromInput: true
    - name: name
    - name: colour
    - name: food:[String]
    - name: legs:Int
  
    Food:
    - name: id:ID!
    - name: name
    - name: tasty:Bool
`)

	type Field struct {
		Name             string
		ExcludeFromInput bool `yaml:"excludeFromInput,omitempty"`
	}

	type Object []Field

	type Schema struct {
		Objects map[string]Object
	}

	var s Schema

	err := yaml.Unmarshal(y, &s)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(s)
}
