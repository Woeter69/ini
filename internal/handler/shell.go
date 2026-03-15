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
	Register("shell", &ShellHandler{})
}

// ShellHandler scaffolds shell script projects.
type ShellHandler struct{}

func (s *ShellHandler) Name() string {
	return "Shell"
}

// SupportedTypes declares which global taxonomy categories Shell supports
func (s *ShellHandler) SupportedTypes() []string {
	return []string{"basic", "devops", "network", "os", "security", "script"}
}

func (s *ShellHandler) Validate() error {
	return nil
}

func (s *ShellHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory + lib/
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	os.MkdirAll(filepath.Join(projectDir, "lib"), 0o755)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("shell/%s/main.sh.tmpl", typeDir)
	utilsTmplPath := "shell/lib/utils.sh.tmpl"

	// 2. Create main script
	scriptName := config.Name + ".sh"
	fmt.Printf("  %s Creating %s...\n", ui.Arrow, scriptName)
	if err := s.processTemplate(config, mainTmplPath, filepath.Join(projectDir, scriptName), 0o755); err != nil {
		return err
	}
	fmt.Printf("  %s %s created (executable)\n", ui.CheckMark, scriptName)

	// 3. Create lib/utils.sh
	fmt.Printf("  %s Creating lib/utils.sh...\n", ui.Arrow)
	if err := s.processTemplate(config, utilsTmplPath, filepath.Join(projectDir, "lib", "utils.sh"), 0o755); err != nil {
		return err
	}
	fmt.Printf("  %s lib/utils.sh created\n", ui.CheckMark)

	// 4. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 5. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 6. Initialize git repo (if --git flag is set)
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	// Print summary
	fmt.Println()
	relPath, _ := filepath.Rel(".", projectDir)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Shell %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  ./%s.sh\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}

func (s *ShellHandler) processTemplate(config ProjectConfig, tmplPath, destPath string, perm os.FileMode) error {
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

	if err := os.WriteFile(destPath, buf.Bytes(), perm); err != nil {
		return fmt.Errorf("failed to write %s: %w", destPath, err)
	}
	return nil
}
