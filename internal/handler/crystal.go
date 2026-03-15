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

func init() { Register("crystal", &CrystalHandler{}) }

type CrystalHandler struct{}

func (c *CrystalHandler) Name() string { return "Crystal" }

// SupportedTypes declares which global taxonomy categories Crystal supports
func (c *CrystalHandler) SupportedTypes() []string {
	return []string{"basic", "web", "cli", "data"}
}

func (c *CrystalHandler) Validate() error {
	if _, err := exec.LookPath("crystal"); err != nil {
		return fmt.Errorf("crystal is not installed or not in PATH.\n  Install it: https://crystal-lang.org/install/")
	}
	return nil
}

func (c *CrystalHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	fmt.Printf("  %s Initializing Crystal project...\n", ui.Arrow)
	parent := filepath.Dir(projectDir)
	if parent == "" {
		parent = "."
	}
	crystalInit := exec.Command("crystal", "init", "app", config.Name)
	crystalInit.Dir = parent
	crystalInit.Stdout = nil
	crystalInit.Stderr = os.Stderr
	if err := crystalInit.Run(); err != nil {
		return fmt.Errorf("failed to init crystal project: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	shardTmplPath := fmt.Sprintf("crystal/%s/shard.yml.tmpl", typeDir)
	mainTmplPath := fmt.Sprintf("crystal/%s/src/project.cr.tmpl", typeDir)

	// 1. Overwrite shard.yml
	fmt.Printf("  %s Customizing shard.yml...\n", ui.Arrow)
	if err := c.processTemplate(config, shardTmplPath, filepath.Join(projectDir, "shard.yml")); err != nil {
		return err
	}

	// 2. Overwrite src/{name}.cr
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	if err := c.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "src", config.Name+".cr")); err != nil {
		return err
	}
	fmt.Printf("  %s Crystal project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		os.RemoveAll(filepath.Join(projectDir, ".git")) // crystal init app creates git
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	c.printSummary(config)
	return nil
}

func (c *CrystalHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (c *CrystalHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Crystal %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  crystal run src/%s.cr\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
