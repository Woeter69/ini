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

func init() { Register("ruby", &RubyHandler{}) }

type RubyHandler struct{}

func (r *RubyHandler) Name() string { return "Ruby" }

// SupportedTypes declares which global taxonomy categories Ruby supports
func (r *RubyHandler) SupportedTypes() []string {
	return []string{"basic", "web", "cli", "gem"}
}

func (r *RubyHandler) Validate() error { return nil } // Ruby exists on most systems

func (r *RubyHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}

	// 1. Create main.rb (or similar file)
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	mainTmplPath := fmt.Sprintf("ruby/%s/main.rb.tmpl", typeDir)
	mainDestPath := filepath.Join(projectDir, "main.rb")
	if typeDir == "gem" {
		mainDestPath = filepath.Join(projectDir, "lib", config.Name+".rb")
		if err := os.MkdirAll(filepath.Join(projectDir, "lib"), 0o755); err != nil {
			return err
		}
		// Gemspec is unique
		if err := r.processTemplate(config, "ruby/gem/project.gemspec.tmpl", filepath.Join(projectDir, config.Name+".gemspec")); err != nil {
			return err
		}
	} else {
		if err := r.processTemplate(config, mainTmplPath, mainDestPath); err != nil {
			return err
		}
	}

	// 2. Create Gemfile (if applicable)
	if typeDir != "gem" {
		fmt.Printf("  %s Creating Gemfile...\n", ui.Arrow)
		gemfileTmplPath := fmt.Sprintf("ruby/%s/Gemfile.tmpl", typeDir)
		if err := r.processTemplate(config, gemfileTmplPath, filepath.Join(projectDir, "Gemfile")); err != nil {
			return err
		}
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

	r.printSummary(config)
	return nil
}

func (r *RubyHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (r *RubyHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Ruby %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  ruby main.rb\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
