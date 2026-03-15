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

func init() { Register("nim", &NimHandler{}) }

type NimHandler struct{}

func (n *NimHandler) Name() string { return "Nim" }

// SupportedTypes declares which global taxonomy categories Nim supports
func (n *NimHandler) SupportedTypes() []string {
	return []string{"basic", "app", "api", "cli", "web", "data", "game", "embedded", "math"}
}

func (n *NimHandler) Validate() error {
	if _, err := exec.LookPath("nimble"); err != nil {
		return fmt.Errorf("nimble is not installed or not in PATH.\n  Install it: https://nim-lang.org/install.html")
	}
	return nil
}

func (n *NimHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Nim project with nimble...\n", ui.Arrow)
	nimbleInit := exec.Command("nimble", "init", "-y")
	nimbleInit.Dir = projectDir
	if err := nimbleInit.Run(); err != nil {
		return fmt.Errorf("failed to init nimble project: %w", err)
	}
	fmt.Printf("  %s nimble project structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "api" {
		typeDir = "basic"
	}
	if typeDir == "db" {
		typeDir = "data"
	}
	mainTmplPath := fmt.Sprintf("nim/%s/main.nim.tmpl", typeDir)
	nimbleTmplPath := fmt.Sprintf("nim/%s/project.nimble.tmpl", typeDir)

	// 3. Overwrite src/{name}.nim
	fmt.Printf("  %s Generating src/%s.nim...\n", ui.Arrow, config.Name)
	mainPath := filepath.Join(projectDir, "src", config.Name+".nim")
	if err := n.processTemplate(config, mainTmplPath, mainPath); err != nil {
		return err
	}

	// 4. Overwrite {name}.nimble
	fmt.Printf("  %s Configuring %s.nimble...\n", ui.Arrow, config.Name)
	nimblePath := filepath.Join(projectDir, config.Name+".nimble")
	if err := n.processTemplate(config, nimbleTmplPath, nimblePath); err != nil {
		// Fallback to basic
	}

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

	n.printSummary(config)
	return nil
}

func (n *NimHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback for nimble
		if strings.HasSuffix(tmplPath, "project.nimble.tmpl") {
			tmplPath = "nim/basic/project.nimble.tmpl"
			content, err = templates.FS.ReadFile(tmplPath)
		}
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
		}
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

func (n *NimHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Nim %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  nim c -r src/%s.nim\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
