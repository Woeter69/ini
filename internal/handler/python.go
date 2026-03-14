package handler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
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

	// 2. Create virtual environment with uv
	fmt.Printf("  %s Creating virtual environment with uv...\n", ui.Arrow)
	uvVenv := exec.Command("uv", "venv")
	uvVenv.Dir = projectDir
	uvVenv.Stdout = os.Stdout
	uvVenv.Stderr = os.Stderr
	if err := uvVenv.Run(); err != nil {
		return fmt.Errorf("failed to create venv: %w", err)
	}
	fmt.Printf("  %s Virtual environment created\n", ui.CheckMark)

	// 3. Create requirements.txt
	fmt.Printf("  %s Creating requirements.txt...\n", ui.Arrow)
	reqContent := "# Add your project dependencies here\n# Example:\n# requests>=2.31.0\n# flask>=3.0.0\n"
	if err := os.WriteFile(filepath.Join(projectDir, "requirements.txt"), []byte(reqContent), 0o644); err != nil {
		return fmt.Errorf("failed to create requirements.txt: %w", err)
	}
	fmt.Printf("  %s requirements.txt created\n", ui.CheckMark)

	// 4. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 5. Create main.py
	fmt.Printf("  %s Creating main.py...\n", ui.Arrow)
	mainPy := fmt.Sprintf(`"""
%s — A Python project.
"""


def main():
    print("Hello from %s!")


if __name__ == "__main__":
    main()
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(projectDir, "main.py"), []byte(mainPy), 0o644); err != nil {
		return fmt.Errorf("failed to create main.py: %w", err)
	}
	fmt.Printf("  %s main.py created\n", ui.CheckMark)

	// 6. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 7. Initialize git repo (if --git flag is set)
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
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Python project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  source .venv/bin/activate\n")
	summary.WriteString("  python main.py\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
