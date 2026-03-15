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

func init() { Register("julia", &JuliaHandler{}) }

type JuliaHandler struct{}

func (j *JuliaHandler) Name() string { return "Julia" }

// SupportedTypes declares which global taxonomy categories Julia supports
func (j *JuliaHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "math", "data", "stats"}
}

func (j *JuliaHandler) Validate() error { return nil }

func (j *JuliaHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" {
		typeDir = "basic"
	}
	if typeDir == "cli" {
		typeDir = "script"
	}

	fmt.Printf("  %s Generating main.jl...\n", ui.Arrow)
	tmplPath := fmt.Sprintf("julia/%s/main.jl.tmpl", typeDir)
	if err := j.processTemplate(config, tmplPath, filepath.Join(projectDir, "main.jl")); err != nil {
		return err
	}

	fmt.Printf("  %s Julia project initialized\n", ui.CheckMark)

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

	j.printSummary(config)
	return nil
}

func (j *JuliaHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (j *JuliaHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Julia %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  julia main.jl\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
