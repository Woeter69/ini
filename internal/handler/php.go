package handler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/templates"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("php", &PHPHandler{}) }

type PHPHandler struct{}

func (p *PHPHandler) Name() string { return "PHP" }

// SupportedTypes declares which global taxonomy categories PHP supports
func (p *PHPHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "web", "api"}
}

func (p *PHPHandler) Validate() error {
	if _, err := exec.LookPath("php"); err != nil {
		return fmt.Errorf("php is not installed or not in PATH.\n  Install it: https://www.php.net/downloads")
	}
	return nil
}

func (p *PHPHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create src/ directory
	os.MkdirAll(filepath.Join(projectDir, "src"), 0o755)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("php/%s/main.php.tmpl", typeDir)
	composerTmplPath := fmt.Sprintf("php/%s/composer.json.tmpl", typeDir)

	// 1. Create main.php
	fmt.Printf("  %s Generating main.php...\n", ui.Arrow)
	if err := p.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "main.php")); err != nil {
		return err
	}

	// 2. Create composer.json
	fmt.Printf("  %s Creating composer.json...\n", ui.Arrow)
	if err := p.processTemplate(config, composerTmplPath, filepath.Join(projectDir, "composer.json")); err != nil {
		return err
	}
	fmt.Printf("  %s PHP project initialized\n", ui.CheckMark)

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

	p.printSummary(config)
	return nil
}

func (p *PHPHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (p *PHPHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your PHP %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  php main.php\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
