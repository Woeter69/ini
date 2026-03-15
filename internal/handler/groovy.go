package handler

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/templates"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("groovy", &GroovyHandler{}) }

type GroovyHandler struct{}

func (g *GroovyHandler) Name() string { return "Groovy" }

// SupportedTypes declares which global taxonomy categories Groovy supports
func (g *GroovyHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "api", "cli", "data"}
}

func (g *GroovyHandler) Validate() error { return nil }

func (g *GroovyHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	if typeDir == "api" {
		typeDir = "web"
	}
	if typeDir == "data" {
		typeDir = "app" // Use Gradle app for data if specific template missing
	}

	switch typeDir {
	case "basic", "cli":
		fmt.Printf("  %s Generating main.groovy...\n", ui.Arrow)
		tmplPath := fmt.Sprintf("groovy/%s/main.groovy.tmpl", typeDir)
		if err := g.processTemplate(config, tmplPath, filepath.Join(projectDir, "main.groovy")); err != nil {
			return err
		}
	case "web":
		fmt.Printf("  %s Generating Ratpack project...\n", ui.Arrow)
		// build.gradle
		if err := g.processTemplate(config, "groovy/web/build.gradle.tmpl", filepath.Join(projectDir, "build.gradle")); err != nil {
			return err
		}
		// src/ratpack/ratpack.groovy
		os.MkdirAll(filepath.Join(projectDir, "src", "ratpack"), 0o755)
		if err := g.processTemplate(config, "groovy/web/src/ratpack/ratpack.groovy.tmpl", filepath.Join(projectDir, "src", "ratpack", "ratpack.groovy")); err != nil {
			return err
		}
	case "app":
		fmt.Printf("  %s Generating structured Groovy app...\n", ui.Arrow)
		// build.gradle
		if err := g.processTemplate(config, "groovy/app/build.gradle.tmpl", filepath.Join(projectDir, "build.gradle")); err != nil {
			return err
		}
		// src/main/groovy/{name}/App.groovy
		packageDir := filepath.Join(projectDir, "src", "main", "groovy", config.Name)
		os.MkdirAll(packageDir, 0o755)
		if err := g.processTemplate(config, "groovy/app/src/main/groovy/project/App.groovy.tmpl", filepath.Join(packageDir, "App.groovy")); err != nil {
			return err
		}
	}

	fmt.Printf("  %s Groovy project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	g.printSummary(config)
	return nil
}

func (g *GroovyHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	t, err := template.New(filepath.Base(tmplPath)).Delims("[[", "]]").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(destPath, buf.Bytes(), 0o644)
}

func (g *GroovyHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Groovy %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	if config.Type == "basic" || config.Type == "cli" {
		summary.WriteString("  groovy main.groovy\n")
	} else {
		summary.WriteString("  gradle run\n")
	}

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
