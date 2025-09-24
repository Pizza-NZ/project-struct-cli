package builder

import (
	"io"
	"pizza-nz/project-struct-cli/templates"
)

// LLMBuilder is a concrete implementation of DocumentBuilder that generates
// a Markdown document from the provided file data.
type LLMBuilder struct {
	projectName    string
	projectSummary templates.FileData
	files          []templates.FileData
}

// NewLLMBuilder creates and returns a new LLMBuilder instance.
func NewLLMBuilder() *LLMBuilder {
	return &LLMBuilder{}
}

// SetProjectName stores the project name for use in the template.
func (b *LLMBuilder) SetProjectName(name string) {
	b.projectName = name
}

// // SetFileTree is a placeholder for a future feature.
// func (b *LLMBuilder) SetFileTree(tree string) {
// 	// TODO: Implement file tree generation and inclusion.
// }

func (b *LLMBuilder) SetSummary(summary templates.FileData) {
	b.projectSummary = summary
}

// AddFile appends file data to the internal slice.
func (b *LLMBuilder) AddFile(file templates.FileData) {
	b.files = append(b.files, file)
}

// Build generates the final document with template.
// It returns the generated document as an io.Reader.
func (b *LLMBuilder) Build() (io.Reader, error) {
	templateData := templates.TemplateData{
		ProjectName:    b.projectName,
		ProjectSummary: b.projectSummary,
		Files:          b.files,
	}
	return templates.ExecuteTemplate(templates.LLM, templateData)
}
