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

func init() { Register("haskell", &HaskellHandler{}) }

type HaskellHandler struct{}

func (h *HaskellHandler) Name() string { return "Haskell" }

// SupportedTypes declares which global taxonomy categories Haskell supports
func (h *HaskellHandler) SupportedTypes() []string {
	return []string{"basic", "cli", "web", "ai"}
}

func (h *HaskellHandler) Validate() error {
	if _, err := exec.LookPath("cabal"); err != nil {
		return fmt.Errorf("cabal is not installed or not in PATH.\n  Install it: https://www.haskell.org/ghcup/")
	}
	return nil
}

func (hs *HaskellHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Haskell project with cabal...\n", ui.Arrow)
	cabalInit := exec.Command("cabal", "init", "--non-interactive", "--package-name", config.Name)
	cabalInit.Dir = projectDir
	if err := cabalInit.Run(); err != nil {
		return fmt.Errorf("failed to init cabal project: %w", err)
	}
	fmt.Printf("  %s Cabal project structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	cabalTmplPath := fmt.Sprintf("haskell/%s/project.cabal.tmpl", typeDir)
	mainTmplPath := fmt.Sprintf("haskell/%s/Main.hs.tmpl", typeDir)

	// 1. Overwrite .cabal file
	fmt.Printf("  %s Configuring .cabal file...\n", ui.Arrow)
	cabalFile := config.Name + ".cabal"
	if err := hs.processTemplate(config, cabalTmplPath, filepath.Join(projectDir, cabalFile)); err != nil {
		return err
	}

	// 2. Overwrite app/Main.hs
	fmt.Printf("  %s Generating app/Main.hs...\n", ui.Arrow)
	if err := os.MkdirAll(filepath.Join(projectDir, "app"), 0o755); err != nil {
		return err
	}
	if err := hs.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "app", "Main.hs")); err != nil {
		return err
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

	hs.printSummary(config)
	return nil
}

func (hs *HaskellHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (hs *HaskellHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Haskell %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  cabal run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
