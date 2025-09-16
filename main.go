package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	ignore "github.com/sabhiram/go-gitignore"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type Config struct {
	SrcDir     string
	OutputFile string
	IgnoreStr  string
	Format     string
}

type Generator struct {
	ignoreMatcher ignore.IgnoreParser

	SrcDir string
	Files  []FileData
}

type FileData struct {
	Path     string
	Content  string
	Language string
}

type TemplateData struct {
	ProjectName string
	Files       []FileData
}

func NewGenerator(matcher ignore.IgnoreParser, src string) *Generator {
	return &Generator{
		ignoreMatcher: matcher,
		SrcDir:        src,
		Files:         make([]FileData, 0),
	}
}

func (g *Generator) processPath(path string, d os.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if g.ignoreMatcher != nil && g.ignoreMatcher.MatchesPath(path) {
		if d.IsDir() {
			return filepath.SkipDir
		}

		return nil
	}

	if !d.IsDir() {
		fmt.Printf("Processing file: %s\n", path)

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			log.Printf("Could not read file %s: %v", path, readErr)
			return nil
		}

		relativePath, err := filepath.Rel(g.SrcDir, path)
		if err != nil {
			// Handle error, but you can probably just use the original path
			relativePath = path
		}

		file := FileData{
			Path:     relativePath,
			Content:  string(content),
			Language: getFileLanguage(path),
		}
		g.Files = append(g.Files, file)
	}

	return nil
}

func run(cfg Config, output io.Writer) error {
	fmt.Printf("Scanning directory: %s\n", cfg.SrcDir)

	var ignoreMatcher ignore.IgnoreParser
	gitignorePath := filepath.Join(cfg.SrcDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		matcher, err := ignore.CompileIgnoreFile(gitignorePath)
		if err != nil {
			log.Printf("Could not parse .gitignore file: %v", err)
		} else {
			ignoreMatcher = matcher
			fmt.Println("Loaded .gitignore rules.")
		}
	}

	generator := NewGenerator(ignoreMatcher, cfg.SrcDir)

	if err := filepath.WalkDir(cfg.SrcDir, generator.processPath); err != nil {
		return err
	}

	templateData := TemplateData{
		ProjectName: filepath.Base(cfg.SrcDir),
		Files:       generator.Files,
	}

	var templatePath string
	switch cfg.Format {
	case "review":
		templatePath = "templates/review.md.tmpl"
	case "llm":
		templatePath = "templates/llm.txt.tmpl"
	case "default":
		templatePath = "templates/default.md.tmpl"
	default:
		return fmt.Errorf("unknown format: %s", cfg.Format)
	}

	templateBytes, err := templatesFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read embedded template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("default").Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	return tmpl.Execute(output, templateData)
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.SrcDir, "src", ".", "The source directory to scan.")
	flag.StringVar(&cfg.OutputFile, "out", "project_structure.md", "The name of the output document.")
	flag.StringVar(&cfg.IgnoreStr, "ignore", ".git,.idea,node_modules,vendor,build,dist", "Legacy: still here but .gitignore is preferred.")
	flag.StringVar(&cfg.Format, "format", "default", "The output format for the document (e.g., default, review, llm).")
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
		return ""
	}
}
