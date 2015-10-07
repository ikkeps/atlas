package main

import (
	"fmt"
	"text/template"
)

// Available descriptor templates
const (
	DESC_KIWI DescriptorFormat = "kiwi"
)

// Represents a Descriptor Format
type DescriptorFormat string

// Get the template for the given descriptor format
// Returns an error if the template can not be parsed
func GetTemplateForFormat(format DescriptorFormat) (*template.Template, error) {
	t := template.New(fmt.Sprintf("%s.template", format))
	return t.ParseFiles(fmt.Sprintf("templates/%s.template", format))
}
