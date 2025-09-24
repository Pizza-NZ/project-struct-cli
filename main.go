// project-struct-cli is a command-line tool that scans a directory and
// consolidates all of its file contents into a single, structured document.
//
// It is designed to help create comprehensive project contexts for code reviews,
// documentation, or for use with Large Language Models (LLMs). The tool is
// configurable via command-line flags.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	builder "pizza-nz/project-struct-cli/builders"
	"pizza-nz/project-struct-cli/templates"
)

// Config holds all the configuration parameters for the application,
// primarily gathered from command-line flags.
type Config struct {
	// SrcDir is the source directory that will be scanned.
	SrcDir string
	// OutputFile is the path to the file where the final document will be written.
	OutputFile string
	// IgnoreCli is a legacy comma-separated string of patterns to ignore.
	// Using a .gitignore file is the preferred method.
	IgnoreCli string
	// Format specifies which output template to use (e.g., 'default', 'llm').
	Format string
	// MaxSizeKB specifies the maximum individual file size in KB to include.
	MaxSizeKB int64
	// MaxTotalSizeMB specifies the maximum total size of all files in MB.
	MaxTotalSizeMB int64
}

// --- Main Application Logic ---

// run is the main logic function of the application. It orchestrates the
// process of building the document from the given configuration.
func run(cfg Config, output io.Writer) error {
	var build builder.DocumentBuilder
	switch cfg.Format {
	case "review":
		build = builder.NewReviewBuilder()
	case "llm":
		build = builder.NewLLMBuilder()
	default: // "default"
		build = builder.NewDefaultBuilder()
	}
	build.SetProjectName(filepath.Base(cfg.SrcDir))

	readmePath := filepath.Join(cfg.SrcDir, "README.md")
	if content, err := os.ReadFile(readmePath); err == nil {
		file := templates.FileData{
			Path:     "",
			Content:  string(content),
			Language: getFileLanguage(readmePath),
		}
		build.SetSummary(file)
	}

	// Set up the generator with the builder, source directory, and .gitignore file.
	generator := NewGenerator(
		WithBuilder(build),
		WithSrcDir(cfg.SrcDir),
		WithCliIgnore(cfg.IgnoreCli),
		WithMaxSize(cfg.MaxSizeKB),
		WithTotalSizeLimit(cfg.MaxTotalSizeMB),
	)

	// Walk the directory tree and collect file data.
	if err := generator.Walk(); err != nil {
		return err
	}

	// Build the final document using the specified format.
	doc, err := build.Build()
	if err != nil {
		return err
	}

	// Write the generated document to the output.
	_, err = io.Copy(output, doc)
	return err
}

// main is the application entry point. It parses command-line flags,
// sets up the output file, and calls the run function to execute the logic.
func main() {
	cfg := Config{}
	flag.StringVar(&cfg.SrcDir, "src", ".", "The source directory to scan.")
	flag.StringVar(&cfg.OutputFile, "out", "project_structure.md", "The name of the output document.")
	flag.StringVar(&cfg.IgnoreCli, "ignore", ".idea,node_modules,vendor,build,dist", "Comma-separated list of file patterns to ignore.")
	flag.StringVar(&cfg.Format, "format", "default", "The output format for the document (e.g., default, review, llm).")
	flag.Int64Var(&cfg.MaxSizeKB, "max-size", 2048, "Maximum individual file size in KB to include (e.g., 2048 for 2MB).")
	flag.Int64Var(&cfg.MaxTotalSizeMB, "max-total", 100, "Maximum total size of all files in MB.")
	flag.Parse()

	f, err := os.Create(cfg.OutputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	if err := run(cfg, f); err != nil {
		log.Fatalf("Error during directory scan: %v", err)
	}

	fmt.Printf("\nSuccess! Project structure written to %s\n", cfg.OutputFile)
}

// --- Helper Functions ---

// getFileLanguage determines a file's programming language based on its extension.
// It returns a string suitable for use in Markdown code blocks for syntax highlighting.
func getFileLanguage(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".yaml", ".yml":
		return "yaml"
	case ".sh":
		return "shell"
	default:
		// Return an empty string for unknown types, so Markdown will not
		// try to apply syntax highlighting.
		return ""
	}
}
