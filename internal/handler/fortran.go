package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
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
	parentDir := projectDir[:len(projectDir)-len(config.Name)-1]
	if parentDir == "" {
		parentDir = "."
	}
	fpmNew.Dir = parentDir
	fpmNew.Stdout = nil
	fpmNew.Stderr = os.Stderr
	if err := fpmNew.Run(); err != nil {
		return fmt.Errorf("failed to create fpm project: %w", err)
	}
	fmt.Printf("  %s fpm project created\n", ui.CheckMark)

	// 2. Overwrite .gitignore with our comprehensive one (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 3. Overwrite README.md with our template (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 4. Initialize git repo (if --git flag is set)
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	// Print summary
	fmt.Println()
	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Fortran project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	summary.WriteString("  fpm run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
