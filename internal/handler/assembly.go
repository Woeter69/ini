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

func init() {
	Register("assembly", &AssemblyHandler{})
}

// AssemblyHandler scaffolds x86-64 Assembly projects with NASM and a Makefile.
type AssemblyHandler struct{}

func (a *AssemblyHandler) Name() string {
	return "Assembly"
}

func (a *AssemblyHandler) Validate() error {
	return nil
}

// SupportedTypes declares which global taxonomy categories Assembly supports
func (a *AssemblyHandler) SupportedTypes() []string {
	return []string{"basic", "os", "embedded", "cli"}
}

func (a *AssemblyHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	srcDir := filepath.Join(projectDir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		return fmt.Errorf("failed to create src/: %w", err)
	}

	// Determine template paths
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "cli" {
		typeDir = "basic"
	}

	mainFile := "main.asm"
	if typeDir == "os" {
		mainFile = "boot.asm"
	}

	mainTmplPath := fmt.Sprintf("assembly/%s/%s.tmpl", typeDir, mainFile)
	makefileTmplPath := fmt.Sprintf("assembly/%s/Makefile.tmpl", typeDir)

	fmt.Printf("  %s Generating Assembly %s boilerplate...\n", ui.Arrow, typeDir)

	if err := a.processTemplate(config, mainTmplPath, filepath.Join(srcDir, mainFile)); err != nil {
		return err
	}
	if err := a.processTemplate(config, makefileTmplPath, filepath.Join(projectDir, "Makefile")); err != nil {
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

	a.printSummary(config)
	return nil
}

func (a *AssemblyHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback to basic
		if strings.Contains(tmplPath, "main.asm") || strings.Contains(tmplPath, "boot.asm") {
			tmplPath = "assembly/basic/main.asm.tmpl"
			content, err = templates.FS.ReadFile(tmplPath)
		} else if strings.Contains(tmplPath, "Makefile") {
			tmplPath = "assembly/basic/Makefile.tmpl"
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

func (a *AssemblyHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Assembly %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	if config.Type == "os" {
		summary.WriteString("  make && make run (requires qemu)\n")
	} else {
		summary.WriteString("  make && make run\n")
	}

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
