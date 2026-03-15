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

func init() { Register("ocaml", &OCamlHandler{}) }

type OCamlHandler struct{}

func (o *OCamlHandler) Name() string { return "OCaml" }
func (o *OCamlHandler) Validate() error {
	if _, err := exec.LookPath("dune"); err != nil {
		return fmt.Errorf("dune is not installed or not in PATH.\n  Install it: opam install dune")
	}
	return nil
}

// SupportedTypes declares which global taxonomy categories OCaml supports
func (o *OCamlHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "web"}
}

func (o *OCamlHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing OCaml project with dune...\n", ui.Arrow)
	// We use dune init project to get the library/test structure, then overwrite
	duneInit := exec.Command("dune", "init", "project", config.Name, "--non-interactive")
	parent := filepath.Dir(projectDir)
	if parent == "" {
		parent = "."
	}
	duneInit.Dir = parent
	if err := duneInit.Run(); err != nil {
		// Fallback: manually create some dirs if dune init fails in mock environments
		os.MkdirAll(filepath.Join(projectDir, "bin"), 0o755)
		os.MkdirAll(filepath.Join(projectDir, "lib"), 0o755)
		os.MkdirAll(filepath.Join(projectDir, "test"), 0o755)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("ocaml/%s/bin/main.ml.tmpl", typeDir)
	duneTmplPath := "ocaml/basic/bin/dune.tmpl"
	projectTmplPath := "ocaml/basic/dune-project.tmpl"

	// 2. Overwrite files from templates
	fmt.Printf("  %s Generating boilerplate...\n", ui.Arrow)
	if err := o.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "bin", "main.ml")); err != nil {
		return err
	}
	if err := o.processTemplate(config, duneTmplPath, filepath.Join(projectDir, "bin", "dune")); err != nil {
		return err
	}
	if err := o.processTemplate(config, projectTmplPath, filepath.Join(projectDir, "dune-project")); err != nil {
		return err
	}

	fmt.Printf("  %s OCaml project initialized\n", ui.CheckMark)

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

	o.printSummary(config)
	return nil
}

func (o *OCamlHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback to basic main.ml if type-specific one is missing
		if strings.HasSuffix(tmplPath, "main.ml.tmpl") {
			tmplPath = "ocaml/basic/bin/main.ml.tmpl"
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

func (o *OCamlHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your OCaml %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  dune build\n")
	summary.WriteString(fmt.Sprintf("  dune exec ./%s\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
