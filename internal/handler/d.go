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

func init() { Register("d", &DHandler{}) }

type DHandler struct{}

func (d *DHandler) Name() string { return "D" }

// SupportedTypes declares which global taxonomy categories D supports
func (d *DHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "web", "api", "game", "data"}
}

func (d *DHandler) Validate() error {
	if _, err := exec.LookPath("dub"); err != nil {
		return fmt.Errorf("dub is not installed or not in PATH.\n  Install it: https://dlang.org/download.html")
	}
	return nil
}

func (d *DHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing D project with dub...\n", ui.Arrow)
	dubInit := exec.Command("dub", "init", "-n") // -n non-interactive
	dubInit.Dir = projectDir
	dubInit.Stdout = nil
	dubInit.Stderr = os.Stderr
	if err := dubInit.Run(); err != nil {
		return fmt.Errorf("failed to init dub project: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "cli" {
		typeDir = "basic" // Dub handles both via basic structure
	}
	if typeDir == "api" {
		typeDir = "web"
	}
	dubTmplPath := fmt.Sprintf("d/%s/dub.json.tmpl", typeDir)
	appTmplPath := fmt.Sprintf("d/%s/source/app.d.tmpl", typeDir)

	// 1. Overwrite dub.json
	fmt.Printf("  %s Customizing dub.json...\n", ui.Arrow)
	if err := d.processTemplate(config, dubTmplPath, filepath.Join(projectDir, "dub.json")); err != nil {
		return err
	}

	// 2. Overwrite source/app.d
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	os.MkdirAll(filepath.Join(projectDir, "source"), 0o755)
	if err := d.processTemplate(config, appTmplPath, filepath.Join(projectDir, "source", "app.d")); err != nil {
		return err
	}
	fmt.Printf("  %s D project initialized\n", ui.CheckMark)

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

	d.printSummary(config)
	return nil
}

func (d *DHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (d *DHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your D %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  dub run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
