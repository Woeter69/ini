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

func init() { Register("scala", &ScalaHandler{}) }

type ScalaHandler struct{}

func (s *ScalaHandler) Name() string { return "Scala" }

// SupportedTypes declares which global taxonomy categories Scala supports
func (s *ScalaHandler) SupportedTypes() []string {
	return []string{"basic", "cli", "web", "script"}
}

func (s *ScalaHandler) Validate() error {
	if _, err := exec.LookPath("scala-cli"); err != nil {
		return fmt.Errorf("scala-cli is not installed or not in PATH.\n  Install it: https://scala-cli.virtuslab.org/install")
	}
	return nil
}

func (sc *ScalaHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("scala/%s/main.scala.tmpl", typeDir)

	// Create main.scala
	fmt.Printf("  %s Generating main.scala...\n", ui.Arrow)
	if err := sc.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "main.scala")); err != nil {
		return err
	}
	fmt.Printf("  %s Scala project initialized\n", ui.CheckMark)

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

	sc.printSummary(config)
	return nil
}

func (sc *ScalaHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (sc *ScalaHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Scala %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  scala-cli run .\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
