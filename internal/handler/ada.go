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

func init() { Register("ada", &AdaHandler{}) }

type AdaHandler struct{}

func (a *AdaHandler) Name() string { return "Ada" }

// SupportedTypes declares which global taxonomy categories Ada supports
func (a *AdaHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "embedded", "os"}
}

func (a *AdaHandler) Validate() error { return nil }

func (a *AdaHandler) Init(config ProjectConfig) error {
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
	gprTmplPath := "ada/basic/project.gpr.tmpl" // Shared GPR
	mainTmplPath := fmt.Sprintf("ada/%s/src/main.adb.tmpl", typeDir)

	// 1. Create project.gpr
	fmt.Printf("  %s Generating project.gpr...\n", ui.Arrow)
	if err := a.processTemplate(config, gprTmplPath, filepath.Join(projectDir, config.Name+".gpr")); err != nil {
		return err
	}

	// 2. Create src/main.adb
	fmt.Printf("  %s Generating main.adb...\n", ui.Arrow)
	if err := a.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "src", "main.adb")); err != nil {
		return err
	}
	fmt.Printf("  %s Ada project initialized\n", ui.CheckMark)

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

	a.printSummary(config)
	return nil
}

func (a *AdaHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (a *AdaHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Ada %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  gprbuild %s.gpr\n", config.Name))
	summary.WriteString("  ./obj/main\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
