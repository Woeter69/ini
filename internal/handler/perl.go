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

func init() { Register("perl", &PerlHandler{}) }

type PerlHandler struct{}

func (p *PerlHandler) Name() string { return "Perl" }

// SupportedTypes declares which global taxonomy categories Perl supports
func (p *PerlHandler) SupportedTypes() []string {
	return []string{"basic", "cli", "web", "db"}
}

func (p *PerlHandler) Validate() error { return nil }

func (p *PerlHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("perl/%s/main.pl.tmpl", typeDir)
	cpanfileTmplPath := fmt.Sprintf("perl/%s/cpanfile.tmpl", typeDir)

	// 1. Create main.pl
	fmt.Printf("  %s Generating main.pl...\n", ui.Arrow)
	if err := p.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "main.pl")); err != nil {
		return err
	}
	// Make executable
	os.Chmod(filepath.Join(projectDir, "main.pl"), 0o755)

	// 2. Create cpanfile
	fmt.Printf("  %s Creating cpanfile...\n", ui.Arrow)
	if err := p.processTemplate(config, cpanfileTmplPath, filepath.Join(projectDir, "cpanfile")); err != nil {
		return err
	}

	// Create lib/ directory
	os.MkdirAll(filepath.Join(projectDir, "lib"), 0o755)

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

func (p *PerlHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (p *PerlHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Perl %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  perl main.pl\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
