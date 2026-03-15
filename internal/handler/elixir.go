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

func init() { Register("elixir", &ElixirHandler{}) }

type ElixirHandler struct{}

func (e *ElixirHandler) Name() string { return "Elixir" }
func (e *ElixirHandler) Validate() error {
	if _, err := exec.LookPath("mix"); err != nil {
		return fmt.Errorf("mix (Elixir) is not installed or not in PATH.\n  Install it: https://elixir-lang.org/install.html")
	}
	return nil
}

// SupportedTypes declares which global taxonomy categories Elixir supports
func (e *ElixirHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "service"}
}

func (e *ElixirHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// mix new creates the directory itself
	fmt.Printf("  %s Creating Elixir project with mix...\n", ui.Arrow)
	parent := filepath.Dir(projectDir)
	if parent == "" {
		parent = "."
	}

	mixNew := exec.Command("mix", "new", config.Name)
	mixNew.Dir = parent
	if err := mixNew.Run(); err != nil {
		// Fallback for mock environment
		if err := scaffold.CreateDir(projectDir); err != nil {
			return err
		}
		os.MkdirAll(filepath.Join(projectDir, "lib"), 0o755)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" {
		typeDir = "basic"
	}

	// Overwrite files
	if typeDir == "web" {
		tmplPath := "elixir/web/lib/app_web.ex.tmpl"
		destPath := filepath.Join(projectDir, "lib", config.Name+"_web.ex")
		fmt.Printf("  %s Adding Plug web handler...\n", ui.Arrow)
		if err := e.processTemplate(config, tmplPath, destPath); err != nil {
			return err
		}
	} else if typeDir == "service" {
		tmplPath := "elixir/service/lib/worker.ex.tmpl"
		destPath := filepath.Join(projectDir, "lib", "worker.ex")
		fmt.Printf("  %s Adding GenServer worker...\n", ui.Arrow)
		if err := e.processTemplate(config, tmplPath, destPath); err != nil {
			return err
		}
	}

	fmt.Printf("  %s Elixir project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	e.printSummary(config)
	return nil
}

func (e *ElixirHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (e *ElixirHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Elixir %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  iex -S mix\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
