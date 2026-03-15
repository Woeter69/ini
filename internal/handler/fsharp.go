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

func init() { Register("fsharp", &FSharpHandler{}) }

type FSharpHandler struct{}

func (f *FSharpHandler) Name() string { return "F#" }

// SupportedTypes declares which global taxonomy categories F# supports
func (f *FSharpHandler) SupportedTypes() []string {
	return []string{"basic", "web", "db", "ai"}
}

func (f *FSharpHandler) Validate() error {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return fmt.Errorf("dotnet is not installed or not in PATH.\n  Install it: https://dotnet.microsoft.com/download")
	}
	return nil
}

func (f *FSharpHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	fmt.Printf("  %s Creating F# project with dotnet...\n", ui.Arrow)
	dotnetNew := exec.Command("dotnet", "new", "console", "-lang", "F#", "-n", config.Name, "-o", projectDir, "--force")
	dotnetNew.Stdout = nil
	dotnetNew.Stderr = os.Stderr
	if err := dotnetNew.Run(); err != nil {
		return fmt.Errorf("failed to create dotnet F# project: %w", err)
	}
	fmt.Printf("  %s F# project structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	programTmplPath := fmt.Sprintf("fsharp/%s/Program.fs.tmpl", typeDir)
	fsprojTmplPath := fmt.Sprintf("fsharp/%s/project.fsproj.tmpl", typeDir)

	// 1. Overwrite Program.fs
	fmt.Printf("  %s Generating Program.fs...\n", ui.Arrow)
	if err := f.processTemplate(config, programTmplPath, filepath.Join(projectDir, "Program.fs")); err != nil {
		return err
	}

	// 2. Overwrite .fsproj file
	fmt.Printf("  %s Configuring .fsproj file...\n", ui.Arrow)
	fsprojFileName := config.Name + ".fsproj"
	if err := f.processTemplate(config, fsprojTmplPath, filepath.Join(projectDir, fsprojFileName)); err != nil {
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

	f.printSummary(config)
	return nil
}

func (f *FSharpHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (f *FSharpHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your F# %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  dotnet run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
