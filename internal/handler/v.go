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

func init() { Register("v", &VHandler{}) }

type VHandler struct{}

func (v *VHandler) Name() string { return "V" }

// SupportedTypes declares which global taxonomy categories V supports
func (v *VHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "api", "cli", "game", "data"}
}

func (v *VHandler) Validate() error {
	if _, err := exec.LookPath("v"); err != nil {
		return fmt.Errorf("v compiler is not installed or not in PATH.\n  Install it: https://vlang.io/")
	}
	return nil
}

func (v *VHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing V project...\n", ui.Arrow)
	vInit := exec.Command("v", "init")
	vInit.Dir = projectDir
	vInit.Stdout = nil
	vInit.Stderr = os.Stderr
	if err := vInit.Run(); err != nil {
		return fmt.Errorf("failed to init v project: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "cli" {
		typeDir = "basic"
	}
	if typeDir == "api" {
		typeDir = "web"
	}
	modTmplPath := "v/basic/v.mod.tmpl" // Shared v.mod
	mainTmplPath := fmt.Sprintf("v/%s/main.v.tmpl", typeDir)

	// 1. Overwrite v.mod
	fmt.Printf("  %s Customizing v.mod...\n", ui.Arrow)
	if err := v.processTemplate(config, modTmplPath, filepath.Join(projectDir, "v.mod")); err != nil {
		return err
	}

	// 2. Overwrite main.v
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	if err := v.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "main.v")); err != nil {
		return err
	}
	fmt.Printf("  %s V project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		os.RemoveAll(filepath.Join(projectDir, ".git")) // v init creates git
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	v.printSummary(config)
	return nil
}

func (v *VHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (v *VHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your V %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  v run .\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
