package main

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"text/template"
)

// DocumentBuilder defines the interface for constructing a project document.
// This allows for different output builders (e.g., Markdown, HTML) to be used
// interchangeably by the Generator.
type DocumentBuilder interface {
	// SetProjectName sets the name of the project, to be used in document titles.
	SetProjectName(name string)
	// SetFileTree is a planned feature to include a file tree overview in the document.
	SetFileTree(tree string)
	// AddFile adds a file's data to the builder for inclusion in the final document.
	AddFile(file FileData)
	// Build uses a template format to construct the final document and returns it as an io.Reader.
	Build(format string) (io.Reader, error)
}

// MarkdownBuilder is a concrete implementation of DocumentBuilder that generates
// a Markdown document from the provided file data.
type MarkdownBuilder struct {
	projectName string
	files       []FileData
}

// NewMarkdownBuilder creates and returns a new MarkdownBuilder instance.
func NewMarkdownBuilder() *MarkdownBuilder {
	return &MarkdownBuilder{}
}

// SetProjectName stores the project name for use in the template.
func (b *MarkdownBuilder) SetProjectName(name string) {
	b.projectName = name
}

// SetFileTree is a placeholder for a future feature.
func (b *MarkdownBuilder) SetFileTree(tree string) {
	// TODO: Implement file tree generation and inclusion.
}

// AddFile appends file data to the internal slice.
func (b *MarkdownBuilder) AddFile(file FileData) {
	b.files = append(b.files, file)
}

// Build generates the final document by executing a Go template based on the
// provided format string. It returns the generated document as an io.Reader.
func (b *MarkdownBuilder) Build(format string) (io.Reader, error) {
	templateData := TemplateData{
		ProjectName: b.projectName,
		Files:       b.files,
	}

	var templatePath string
	switch format {
	case "review":
		templatePath = "templates/review.md.tmpl"
	case "llm":
		templatePath = "templates/llm.txt.tmpl"
	case "default":
		templatePath = "templates/default.md.tmpl"
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}

	// Read the embedded template file.
	templateBytes, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("could not read embedded template %s: %w", templatePath, err)
	}

	// Parse the template.
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(templateBytes))
	if err != nil {
		return nil, fmt.Errorf("could not parse template: %w", err)
	}

	// Execute the template into a buffer.
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, templateData); err != nil {
		return nil, err
	}

	return &doc, nil
}
