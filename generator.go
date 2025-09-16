package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

// Generator is responsible for walking a directory structure, reading files,
// and passing their data to a DocumentBuilder.
type Generator struct {
	builder       DocumentBuilder
	ignoreMatcher ignore.IgnoreParser
	srcDir        string
}

// Option is a function type used to configure a Generator. This follows the
// "Functional Options" pattern, allowing for flexible and clear configuration.
type Option func(*Generator)

// WithGitIgnore returns an Option that configures the Generator to use a
// .gitignore file for filtering which files and directories to ignore.
func WithGitIgnore(path string) Option {
	return func(g *Generator) {
		matcher, err := ignore.CompileIgnoreFile(path)
		// If the .gitignore file doesn't exist or has errors, we simply
		// proceed without an ignore matcher.
		if err == nil {
			g.ignoreMatcher = matcher
		}
	}
}

// WithBuilder returns an Option that sets the DocumentBuilder for the Generator.
func WithBuilder(builder DocumentBuilder) Option {
	return func(g *Generator) {
		g.builder = builder
	}
}

// WithSrcDir returns an Option that sets the source directory for the Generator to scan.
func WithSrcDir(srcDir string) Option {
	return func(g *Generator) {
		g.srcDir = srcDir
	}
}

// NewGenerator creates a new Generator and applies all the provided functional options.
func NewGenerator(opts ...Option) *Generator {
	g := &Generator{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Walk starts the process of walking the source directory tree.
func (g *Generator) Walk() error {
	fmt.Printf("Scanning directory: %s\n", g.srcDir)
	return filepath.WalkDir(g.srcDir, g.processPath)
}

// processPath is the callback function for filepath.WalkDir. It is called
// for every file and directory in the source tree.
func (g *Generator) processPath(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	// Always skip .git directories.
	if d.IsDir() && d.Name() == ".git" {
		return filepath.SkipDir
	}

	// Check if the path should be ignored based on .gitignore rules.
	if g.ignoreMatcher != nil && g.ignoreMatcher.MatchesPath(path) {
		// If a directory is ignored, skip it entirely.
		if d.IsDir() {
			return filepath.SkipDir
		}
		// If it's an ignored file, just skip this entry.
		return nil
	}

	if !d.IsDir() {
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			// Log the error but don't stop the whole process.
			log.Printf("Could not read file %s: %v", path, readErr)
			return nil
		}

		// Get the file path relative to the source directory for cleaner output.
		relativePath, err := filepath.Rel(g.srcDir, path)
		if err != nil {
			relativePath = path // Fallback to the full path on error.
		}

		file := FileData{
			Path:     relativePath,
			Content:  string(content),
			Language: getFileLanguage(relativePath),
		}
		g.builder.AddFile(file)
	}
	return nil
}
