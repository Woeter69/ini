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

func init() { Register("objc", &ObjCHandler{}) }

type ObjCHandler struct{}

func (o *ObjCHandler) Name() string { return "Objective-C" }

func (o *ObjCHandler) Validate() error { return nil }

// SupportedTypes declares which global taxonomy categories Objective-C supports
func (o *ObjCHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "desktop", "mobile"}
}

func (o *ObjCHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	os.MkdirAll(filepath.Join(projectDir, "src"), 0o755)
	os.MkdirAll(filepath.Join(projectDir, "build"), 0o755)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "cli" {
		typeDir = "basic"
	}

	mainTmplPath := fmt.Sprintf("objc/%s/main.m.tmpl", typeDir)
	makefileTmplPath := "objc/basic/Makefile.tmpl"

	fmt.Printf("  %s Generating Objective-C %s boilerplate...\n", ui.Arrow, typeDir)

	if err := o.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "src", "main.m")); err != nil {
		return err
	}
	if err := o.processTemplate(config, makefileTmplPath, filepath.Join(projectDir, "Makefile")); err != nil {
		return err
	}

	fmt.Printf("  %s Scaffolding complete\n", ui.CheckMark)

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

	o.printSummary(config)
	return nil
}

func (o *ObjCHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback to basic if specific type fails
		if strings.Contains(tmplPath, "main.m") {
			tmplPath = "objc/basic/main.m.tmpl"
			content, err = templates.FS.ReadFile(tmplPath)
		}
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
		}
	}

	t, err := template.New(filepath.Base(tmplPath)).Delims("<<", ">>").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(destPath, buf.Bytes(), 0o644)
}

func (o *ObjCHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Objective-C %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  make\n")
	summary.WriteString(fmt.Sprintf("  ./build/%s\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
