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
	Register("python", &PythonHandler{})
}

// PythonHandler scaffolds Python projects using uv.
type PythonHandler struct{}

func (p *PythonHandler) Name() string {
	return "Python"
}

// SupportedTypes declares which global taxonomy categories Python supports
func (p *PythonHandler) SupportedTypes() []string {
	return []string{"basic", "web", "data", "game", "script"}
}

func (p *PythonHandler) Validate() error {
	_, err := exec.LookPath("uv")
	if err != nil {
		return fmt.Errorf("uv is not installed or not in PATH.\n  Install it: curl -LsSf https://astral.sh/uv/install.sh | sh")
	}
	return nil
}

func (p *PythonHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Initialize uv project (uv init sets up pyproject.toml cleanly)
	fmt.Printf("  %s Initializing Python project with uv...\n", ui.Arrow)
	uvInit := exec.Command("uv", "init", "--app")
	uvInit.Dir = projectDir
	if err := uvInit.Run(); err != nil {
		return fmt.Errorf("failed to init uv project: %w", err)
	}
	fmt.Printf("  %s uv project initialized\n", ui.CheckMark)

	// Determine type path for template
	templatePath := "python/script/main.py.tmpl" // fallback
	deps := []string{}
	cmdStr := "uv run python main.py"

	switch config.Type {
	case "web":
		templatePath = "python/web/main.py.tmpl"
		deps = append(deps, "fastapi", "uvicorn")
	case "data":
		templatePath = "python/data/main.py.tmpl"
		deps = append(deps, "pandas", "numpy")
	case "game":
		templatePath = "python/game/main.py.tmpl"
		deps = append(deps, "pygame")
	case "script", "basic":
		templatePath = "python/script/main.py.tmpl"
	}

	// 3. Add dependencies if any
	if len(deps) > 0 {
		fmt.Printf("  %s Adding dependencies: %s...\n", ui.Arrow, strings.Join(deps, ", "))
		uvAdd := exec.Command("uv", append([]string{"add"}, deps...)...)
		uvAdd.Dir = projectDir
		if err := uvAdd.Run(); err != nil {
			return fmt.Errorf("failed to add dependencies: %w", err)
		}
		fmt.Printf("  %s Dependencies added\n", ui.CheckMark)
	}

	// 4. Create main.py from embedded template
	tmplContent, err := templates.FS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	t, err := template.New("main").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// overwrite the generated hello.py from uv init
	os.Remove(filepath.Join(projectDir, "hello.py"))
	if err := os.WriteFile(filepath.Join(projectDir, "main.py"), buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write main.py: %w", err)
	}
	fmt.Printf("  %s main.py created\n", ui.CheckMark)

	// 5. Create .gitignore and README
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	// 6. Git
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
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Python %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  %s\n", cmdStr))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
