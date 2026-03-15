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

func init() {
	Register("fortran", &FortranHandler{})
}

// FortranHandler scaffolds Fortran projects using fpm.
type FortranHandler struct{}

func (f *FortranHandler) Name() string {
	return "Fortran"
}

// SupportedTypes declares which global taxonomy categories Fortran supports
func (f *FortranHandler) SupportedTypes() []string {
	return []string{"basic", "data", "math", "cli"}
}

func (f *FortranHandler) Validate() error {
	_, err := exec.LookPath("fpm")
	if err != nil {
		return fmt.Errorf("fpm (Fortran Package Manager) is not installed or not in PATH.\n  Install it: https://fpm.fortran-lang.org/install/index.html")
	}
	return nil
}

func (f *FortranHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Run fpm new — it creates the directory and full project structure
	fmt.Printf("  %s Creating Fortran project with fpm...\n", ui.Arrow)
	fpmNew := exec.Command("fpm", "new", config.Name)
	// Run fpm new from the parent directory so it creates the project dir
	parentDir := filepath.Dir(projectDir)
	if parentDir == "" {
		parentDir = "."
	}
	fpmNew.Dir = parentDir
	fpmNew.Stdout = nil
	fpmNew.Stderr = os.Stderr
	if err := fpmNew.Run(); err != nil {
		return fmt.Errorf("failed to create fpm project: %w", err)
	}

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	tomlTmplPath := fmt.Sprintf("fortran/%s/fpm.toml.tmpl", typeDir)
	mainTmplPath := fmt.Sprintf("fortran/%s/app/main.f90.tmpl", typeDir)

	// 2. Overwrite fpm.toml
	fmt.Printf("  %s Customizing fpm.toml...\n", ui.Arrow)
	if err := f.processTemplate(config, tomlTmplPath, filepath.Join(projectDir, "fpm.toml")); err != nil {
		return err
	}

	// 3. Overwrite app/main.f90
	fmt.Printf("  %s Generating source files...\n", ui.Arrow)
	// Ensure app directory exists (fpm new should have created it but let's be safe)
	os.MkdirAll(filepath.Join(projectDir, "app"), 0o755)
	if err := f.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "app", "main.f90")); err != nil {
		return err
	}
	fmt.Printf("  %s Fortran project initialized\n", ui.CheckMark)

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

func (f *FortranHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (f *FortranHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Fortran %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  fpm run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
