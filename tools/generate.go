package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type TemplateGenerator struct {
	TemplatePath string
	Title        string
	Description  string
	Version      string
	Year         string
	Author       string
	WWWPath      string
}

func NewTemplateGenerator(path string, title string, description string, version string, year string, author string, wwwPath string) *TemplateGenerator {
	return &TemplateGenerator{
		TemplatePath: path,
		Title:        title,
		Description:  description,
		Version:      version,
		Year:         year,
		Author:       author,
		WWWPath:      wwwPath,
	}
}

func ReadTemplateFiles(path string) ([]string, error) {
	log.Println("Reading template files from:", path)
	files, err := filepath.Glob(path + "/*.html")
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (t *TemplateGenerator) Generate() {
	files, err := ReadTemplateFiles(t.TemplatePath)
	if err != nil {
		log.Fatalf("Error reading template files: %v", err)
	}

	for _, file := range files {
		log.Println("Generating template:", file)
		tmpl := template.Must(template.ParseFiles(file))
		var output bytes.Buffer
		err := tmpl.Execute(&output, t)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		// Write the generated template to the WWWPath
		outputFile := filepath.Join(t.WWWPath, filepath.Base(file))
		log.Println("Writing template to WWWPath:", outputFile)
		err = os.WriteFile(outputFile, output.Bytes(), 0644)
		if err != nil {
			log.Fatalf("Error writing template to WWWPath: %v", err)
		}

		log.Println("Template written to WWWPath:", outputFile)
	}

	log.Println("Templates generated successfully!")
}

func main() {
	generator := NewTemplateGenerator("templates", "SecNex Reverse Proxy", "SecNex Reverse Proxy", "0.1.0", "2025", "SecNex", "www")
	generator.Generate()
}
