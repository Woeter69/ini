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

func init() { Register("r", &RHandler{}) }

type RHandler struct{}

func (r *RHandler) Name() string { return "R" }

// SupportedTypes declares which global taxonomy categories R supports
func (r *RHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "math", "data", "stats"}
}

func (r *RHandler) Validate() error { return nil }

func (r *RHandler) Init(config ProjectConfig) error {
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

	// 1. Create main.R
	fmt.Printf("  %s Generating main.R...\n", ui.Arrow)
	tmplPath := fmt.Sprintf("r/%s/main.R.tmpl", typeDir)
	if err := r.processTemplate(config, tmplPath, filepath.Join(projectDir, "main.R")); err != nil {
		return err
	}

	// 2. Create R/ directory for functions
	os.MkdirAll(filepath.Join(projectDir, "R"), 0o755)

	// 3. Create .Rprofile
	fmt.Printf("  %s Generating .Rprofile...\n", ui.Arrow)
	rprofile := "# Project-level R settings\n# options(repos = c(CRAN = \"https://cloud.r-project.org\"))\n"
	if err := os.WriteFile(filepath.Join(projectDir, ".Rprofile"), []byte(rprofile), 0o644); err != nil {
		return err
	}

	fmt.Printf("  %s R project initialized\n", ui.CheckMark)

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

	r.printSummary(config)
	return nil
}

func (r *RHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (r *RHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your R %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  Rscript main.R\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
