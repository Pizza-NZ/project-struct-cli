package builder

import (
	"fmt"
	"io"
	"pizza-nz/project-struct-cli/templates"
	"strings"
)

type ReviewBuilder struct {
	projectName    string
	projectSummary templates.FileData
	files          []templates.FileData
}

// ReviewFileData will be used internally by the builder for the review template
type ReviewFileData struct {
	// Path is the relative path of the file from the source directory.
	Path string
	// Language is the detected programming language based on the file extension.
	Language string
	// Lines is the full text content of the file for line numbers.
	Lines []string
}

// NewReviewBuilder creates and returns a new ReviewBuilder instance.
func NewReviewBuilder() *ReviewBuilder {
	return &ReviewBuilder{}
}

// SetProjectName stores the project name for use in the template.
func (b *ReviewBuilder) SetProjectName(name string) {
	b.projectName = name
}

func (b *ReviewBuilder) SetSummary(summary templates.FileData) {
	b.projectSummary = summary
}

// AddFile appends file data to the internal slice.
func (b *ReviewBuilder) AddFile(file templates.FileData) {
	b.files = append(b.files, file)
}

func (b *ReviewBuilder) Build() (io.Reader, error) {
	var reviewFiles []ReviewFileData
	for _, file := range b.files {
		var numberedLines []string
		lines := strings.Split(file.Content, "\n")
		for i, line := range lines {
			// Formated with 4 spaces for number, colon, and line content
			numberedLines = append(numberedLines, fmt.Sprintf("%4d: %s", i+1, line))
		}

		reviewFiles = append(reviewFiles, ReviewFileData{
			Path:     file.Path,
			Language: file.Language,
			Lines:    numberedLines,
		})
	}

	templateData := struct {
		ProjectName    string
		ProjectSummary templates.FileData
		Files          []ReviewFileData
	}{
		ProjectName:    b.projectName,
		ProjectSummary: b.projectSummary,
		Files:          reviewFiles,
	}

	return templates.ExecuteTemplate(templates.Review, templateData)
}
