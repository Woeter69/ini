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

func init() { Register("csharp", &CSharpHandler{}) }

type CSharpHandler struct{}

func (c *CSharpHandler) Name() string { return "C#" }

// SupportedTypes declares which global taxonomy categories C# supports
func (c *CSharpHandler) SupportedTypes() []string {
	return []string{"basic", "web", "db", "desktop", "ai"}
}

func (c *CSharpHandler) Validate() error {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return fmt.Errorf("dotnet is not installed or not in PATH.\n  Install it: https://dotnet.microsoft.com/download")
	}
	return nil
}

func (c *CSharpHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	fmt.Printf("  %s Creating C# project with dotnet...\n", ui.Arrow)
	// We use "console" as a base even for web, since we overwrite the template anyway
	// This ensures we get the basic folder structure and .csproj file
	dotnetNew := exec.Command("dotnet", "new", "console", "-n", config.Name, "-o", projectDir, "--force")
	dotnetNew.Stdout = nil
	dotnetNew.Stderr = os.Stderr
	if err := dotnetNew.Run(); err != nil {
		return fmt.Errorf("failed to create dotnet project: %w", err)
	}
	fmt.Printf("  %s C# project structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	programTmplPath := fmt.Sprintf("csharp/%s/Program.cs.tmpl", typeDir)
	csprojTmplPath := fmt.Sprintf("csharp/%s/project.csproj.tmpl", typeDir)

	// 1. Overwrite Program.cs
	fmt.Printf("  %s Generating Program.cs...\n", ui.Arrow)
	if err := c.processTemplate(config, programTmplPath, filepath.Join(projectDir, "Program.cs")); err != nil {
		return err
	}

	// 2. Overwrite .csproj file
	fmt.Printf("  %s Configuring .csproj file...\n", ui.Arrow)
	csprojFileName := config.Name + ".csproj"
	if err := c.processTemplate(config, csprojTmplPath, filepath.Join(projectDir, csprojFileName)); err != nil {
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

	c.printSummary(config)
	return nil
}

func (c *CSharpHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (c *CSharpHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your C# %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  dotnet run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
