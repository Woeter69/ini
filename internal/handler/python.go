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
	return []string{
		"basic", "web", "desktop", "mobile", "game", "ai", "data", "db", "devops",
		"network", "security", "os", "lang", "finance", "comm", "script", "monitor",
		"stream", "science", "media", "iot", "web3", "graphics", "edu",
	}
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
	case "mobile":
		templatePath = "python/mobile/main.py.tmpl"
		deps = append(deps, "kivy")
	case "desktop":
		templatePath = "python/desktop/main.py.tmpl"
		// tkinter is built-in
	case "game":
		templatePath = "python/game/main.py.tmpl"
		deps = append(deps, "pygame")
	case "ai":
		templatePath = "python/ai/main.py.tmpl"
		deps = append(deps, "torch")
	case "data":
		templatePath = "python/data/main.py.tmpl"
		deps = append(deps, "pandas", "numpy")
	case "db":
		templatePath = "python/db/main.py.tmpl"
		deps = append(deps, "sqlalchemy")
	case "devops":
		templatePath = "python/devops/main.py.tmpl"
		deps = append(deps, "boto3")
	case "network":
		templatePath = "python/network/main.py.tmpl"
		// socket is built-in
	case "security":
		templatePath = "python/security/main.py.tmpl"
		deps = append(deps, "cryptography")
	case "os":
		templatePath = "python/os/main.py.tmpl"
		deps = append(deps, "psutil")
	case "lang":
		templatePath = "python/lang/main.py.tmpl"
		// ast is built-in
	case "finance":
		templatePath = "python/finance/main.py.tmpl"
		deps = append(deps, "yfinance", "pandas")
	case "comm":
		templatePath = "python/comm/main.py.tmpl"
		// smtplib is built-in
	case "monitor":
		templatePath = "python/monitor/main.py.tmpl"
		deps = append(deps, "psutil")
	case "stream":
		templatePath = "python/stream/main.py.tmpl"
		// queue/threading is built-in
	case "science":
		templatePath = "python/science/main.py.tmpl"
		deps = append(deps, "numpy", "scipy")
	case "media":
		templatePath = "python/media/main.py.tmpl"
		deps = append(deps, "pillow")
	case "iot":
		templatePath = "python/iot/main.py.tmpl"
		// random/time built-in
	case "web3":
		templatePath = "python/web3/main.py.tmpl"
		deps = append(deps, "web3")
	case "graphics":
		templatePath = "python/graphics/main.py.tmpl"
		deps = append(deps, "PyOpenGL", "glfw")
	case "edu":
		templatePath = "python/edu/main.py.tmpl"
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
