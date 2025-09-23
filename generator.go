package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	builder "pizza-nz/project-struct-cli/builders"
	"pizza-nz/project-struct-cli/templates"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// A set of common binary and archive file extensions to always ignore.
// Using a map with an empty struct is a memory-efficient way to represent a set.
var binaryExts = map[string]struct{}{
	".exe": {}, ".dll": {}, ".so": {}, ".a": {}, ".lib": {}, ".o": {},
	".zip": {}, ".gz": {}, ".tar": {}, ".rar": {}, ".7z": {},
	".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".bmp": {}, ".ico": {},
	".pdf": {}, ".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {}, ".ppt": {}, ".pptx": {},
}

// Generator is responsible for walking a directory structure, reading files,
// and passing their data to a DocumentBuilder.
type Generator struct {
	builder          builder.DocumentBuilder
	cliIgnoreMatcher ignore.IgnoreParser
	srcDir           string

	matcherCache map[string]ignore.IgnoreParser
}

// Option is a function type used to configure a Generator. This follows the
// "Functional Options" pattern, allowing for flexible and clear configuration.
type Option func(*Generator)

// // WithGitIgnore returns an Option that configures the Generator to use a
// // .gitignore file for filtering which files and directories to ignore.
// func WithGitIgnore(path string) Option {
// 	return func(g *Generator) {
// 		matcher, err := ignore.CompileIgnoreFile(path)
// 		// If the .gitignore file doesn't exist or has errors, we simply
// 		// proceed without an ignore matcher.
// 		if err == nil {
// 			g.gitIgnoreMatcher = matcher
// 		}
// 	}
// }

func WithCliIgnore(patterns string) Option {
	return func(g *Generator) {
		if patterns == "" {
			return // Do nothing if empty
		}

		lines := strings.Split(patterns, ",")
		matcher := ignore.CompileIgnoreLines(lines...)
		g.cliIgnoreMatcher = matcher
	}
}

// WithBuilder returns an Option that sets the DocumentBuilder for the Generator.
func WithBuilder(builder builder.DocumentBuilder) Option {
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
	g := &Generator{
		matcherCache: make(map[string]ignore.IgnoreParser),
	}
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

// findMatcherForDir walks up from a given directory to find the applicalbe .gitingore file.
// It uses a cache to avoid redundant file system lookup.
func (g *Generator) findMatcherForDir(dir string) ignore.IgnoreParser {
	if matcher, exists := g.matcherCache[dir]; exists {
		return matcher
	}

	parentDir := dir

	ignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(ignorePath); err == nil {
		// .gitignore found, compile it
		matcher, err := ignore.CompileIgnoreFile(ignorePath)
		if err == nil {
			// cache the found matcher for the original directory and return
			g.matcherCache[dir] = matcher
			return matcher
		}
	}

	if dir == g.srcDir || dir == "." || dir == "/" {
		g.matcherCache[dir] = nil
		return nil
	}

	parentMatcher := g.findMatcherForDir(parentDir)
	g.matcherCache[dir] = parentMatcher
	return parentMatcher

	// currentDir := dir
	// for {
	// 	ignorePath := filepath.Join(currentDir, ".gitignore")
	// 	if _, err := os.Stat(ignorePath); err == nil {
	// 		// .gitignore found, compile it
	// 		matcher, err := ignore.CompileIgnoreFile(ignorePath)
	// 		if err == nil {
	// 			// cache the found matcher for the original directory and return
	// 			g.matcherCache[dir] = matcher
	// 			return matcher
	// 		}
	// 	}

	// 	if currentDir == g.srcDir || currentDir == "." || currentDir == "/" {
	// 		break
	// 	}

	// 	currentDir = filepath.Dir(currentDir)
	// }

	// g.matcherCache[dir] = nil
	// return nil
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

	dir := filepath.Dir(path)
	gitMatcher := g.findMatcherForDir(dir)

	// Check if the path should be ignored based on .gitignore or cli ignore rules.
	if (gitMatcher != nil && gitMatcher.MatchesPath(path)) ||
		(g.cliIgnoreMatcher != nil && g.cliIgnoreMatcher.MatchesPath(path)) {
		// If a directory is ignored, skip it entirely.
		if d.IsDir() {
			return filepath.SkipDir
		}
		// If it's an ignored file, just skip this entry.
		return nil
	}

	if !d.IsDir() {
		ext := filepath.Ext(path)
		if _, exists := binaryExts[ext]; exists {
			log.Printf("Skipping binary/archive file: %s", path)
			return nil // Skip this file and continue the walk.
		}

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

		file := templates.FileData{
			Path:     relativePath,
			Content:  string(content),
			Language: getFileLanguage(relativePath),
		}
		g.builder.AddFile(file)
	}
	return nil
}
