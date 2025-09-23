package builder

import (
	"io"
	"pizza-nz/project-struct-cli/templates"
)

// DocumentBuilder defines the interface for constructing a project document.
// This allows for different output builders (e.g., Markdown, HTML) to be used
// interchangeably by the Generator.
type DocumentBuilder interface {
	// SetProjectName sets the name of the project, to be used in document titles.
	SetProjectName(name string)
	// // SetFileTree is a planned feature to include a file tree overview in the document.
	// SetFileTree(tree string)
	// SetSummary sets the a README summary
	SetSummary(summary string)
	// AddFile adds a file's data to the builder for inclusion in the final document.
	AddFile(file templates.FileData)
	// Build uses a template format to construct the final document and returns it as an io.Reader.
	Build() (io.Reader, error)
}
