package templates

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

// TemplatesFS holds the embedded template files for generating the documents.
// Using embed allows the binary to be self-contained without needing the
// template files to be present on the filesystem at runtime.
//
//go:embed *.tmpl
var TemplatesFS embed.FS

// TemplatePath defines a type for template file paths to improve type safety.

type TemplatePath string

// Constants for the available template files. Using constants prevents typos
// and makes it easy to manage template paths from one location.
const (
	Review  TemplatePath = "review.md.tmpl"
	Default TemplatePath = "default.md.tmpl"
	LLM     TemplatePath = "llm.txt.tmpl"
)

func (p TemplatePath) String() string {
	return string(p)
}

// TemplateData is the data structure passed to the templates for execution.
type TemplateData struct {
	ProjectName    string
	ProjectSummary string
	Files          []FileData
}

// FileData represents the contents of a single source file.
type FileData struct {
	// Path is the relative path of the file from the source directory.
	Path string
	// Content is the full text content of the file.
	Content string
	// Language is the detected programming language based on the file extension.
	Language string
}

func ExecuteTemplate(templatePath TemplatePath, data any) (io.Reader, error) {
	// Read the embedded template file.
	templateBytes, err := TemplatesFS.ReadFile(templatePath.String())
	if err != nil {
		return nil, fmt.Errorf("could not read embedded template %s: %w", templatePath, err)
	}

	// Parse the template.
	tmpl, err := template.New(filepath.Base(templatePath.String())).Parse(string(templateBytes))
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}

	// Execute the template into a buffer.
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, data); err != nil {
		return nil, err
	}

	return &doc, nil
}
