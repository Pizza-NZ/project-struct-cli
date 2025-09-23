package builder

import (
	"io"
	"pizza-nz/project-struct-cli/templates"
)

// DefaultBuilder is a concrete implementation of DocumentBuilder that generates
// a Markdown document from the provided file data.
type DefaultBuilder struct {
	projectName    string
	projectSummary string
	files          []templates.FileData
}

// NewDefaultBuilder creates and returns a new DefaultBuilder instance.
func NewDefaultBuilder() *DefaultBuilder {
	return &DefaultBuilder{}
}

// SetProjectName stores the project name for use in the template.
func (b *DefaultBuilder) SetProjectName(name string) {
	b.projectName = name
}

// // SetFileTree is a placeholder for a future feature.
// func (b *DefaultBuilder) SetFileTree(tree string) {
// 	// TODO: Implement file tree generation and inclusion.
// }

func (b *DefaultBuilder) SetSummary(summary string) {
	b.projectSummary = summary
}

// AddFile appends file data to the internal slice.
func (b *DefaultBuilder) AddFile(file templates.FileData) {
	b.files = append(b.files, file)
}

// Build generates the final document with template.
// It returns the generated document as an io.Reader.
func (b *DefaultBuilder) Build() (io.Reader, error) {
	templateData := templates.TemplateData{
		ProjectName:    b.projectName,
		ProjectSummary: b.projectSummary,
		Files:          b.files,
	}
	return templates.ExecuteTemplate(templates.Default, templateData)
}
