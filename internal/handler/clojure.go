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

func init() { Register("clojure", &ClojureHandler{}) }

type ClojureHandler struct{}

func (c *ClojureHandler) Name() string { return "Clojure" }

// SupportedTypes declares which global taxonomy categories Clojure supports
func (c *ClojureHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "api", "cli", "data"}
}

func (c *ClojureHandler) Validate() error {
	if _, err := exec.LookPath("clj"); err != nil {
		return fmt.Errorf("clj (Clojure CLI) is not installed.\n  Install it: https://clojure.org/guides/install_clojure")
	}
	return nil
}

func (c *ClojureHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "cli" {
		typeDir = "basic"
	}
	if typeDir == "api" {
		typeDir = "web"
	}
	depsTmplPath := fmt.Sprintf("clojure/%s/deps.edn.tmpl", typeDir)
	coreTmplPath := fmt.Sprintf("clojure/%s/src/project/core.clj.tmpl", typeDir)

	// 1. Create deps.edn
	fmt.Printf("  %s Generating deps.edn...\n", ui.Arrow)
	if err := c.processTemplate(config, depsTmplPath, filepath.Join(projectDir, "deps.edn")); err != nil {
		return err
	}

	// 2. Create src/{name}/core.clj
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	// Sanitize project name for namespace (replaces - with _)
	nsPart := strings.ReplaceAll(config.Name, "-", "_")
	srcDir := filepath.Join(projectDir, "src", nsPart)
	os.MkdirAll(srcDir, 0o755)

	if err := c.processTemplate(config, coreTmplPath, filepath.Join(srcDir, "core.clj")); err != nil {
		return err
	}
	fmt.Printf("  %s Clojure project initialized\n", ui.CheckMark)

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

func (c *ClojureHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	// Namespace needs to be valid clojure symbol part
	config.Name = strings.ReplaceAll(config.Name, "-", "_")

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

func (c *ClojureHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Clojure %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  clj -M:run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
