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

func init() { Register("cobol", &CobolHandler{}) }

type CobolHandler struct{}

func (c *CobolHandler) Name() string { return "COBOL" }

// SupportedTypes declares which global taxonomy categories COBOL supports
func (c *CobolHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "business", "data", "interactive"}
}

func (c *CobolHandler) Validate() error { return nil }

func (c *CobolHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	switch typeDir {
	case "", "basic", "app":
		typeDir = "basic"
	case "cli":
		typeDir = "interactive"
	case "data":
		typeDir = "db"
	}
	mainTmplPath := fmt.Sprintf("cobol/%s/main.cbl.tmpl", typeDir)
	makefileTmplPath := "cobol/basic/Makefile.tmpl" // Shared Makefile

	// 1. Create main.cbl
	fmt.Printf("  %s Generating main.cbl...\n", ui.Arrow)
	if err := c.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "main.cbl")); err != nil {
		return err
	}

	// 2. Create Makefile
	fmt.Printf("  %s Generating Makefile...\n", ui.Arrow)
	if err := c.processTemplate(config, makefileTmplPath, filepath.Join(projectDir, "Makefile")); err != nil {
		return err
	}
	fmt.Printf("  %s COBOL project initialized\n", ui.CheckMark)

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

	c.printSummary(config)
	return nil
}

func (c *CobolHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (c *CobolHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your COBOL %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  make\n")
	summary.WriteString("  ./main\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
