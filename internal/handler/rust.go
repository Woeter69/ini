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
	Register("rust", &RustHandler{})
}

// RustHandler scaffolds Rust projects using cargo.
type RustHandler struct{}

func (r *RustHandler) Name() string {
	return "Rust"
}

func (r *RustHandler) Validate() error {
	_, err := exec.LookPath("cargo")
	if err != nil {
		return fmt.Errorf("cargo is not installed or not in PATH.\n  Install it: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh")
	}
	return nil
}

func (r *RustHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Run cargo init inside the directory
	fmt.Printf("  %s Initializing Rust project with cargo...\n", ui.Arrow)
	cargoInit := exec.Command("cargo", "init", "--name", config.Name)
	cargoInit.Dir = projectDir
	cargoInit.Stdout = nil
	cargoInit.Stderr = os.Stderr
	if err := cargoInit.Run(); err != nil {
		return fmt.Errorf("failed to init cargo project: %w", err)
	}
	fmt.Printf("  %s Cargo project initialized\n", ui.CheckMark)

	// 3. Overwrite .gitignore with our comprehensive one (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 4. Overwrite README.md with our template (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 5. Initialize git repo (if --git flag is set)
	// Note: cargo init creates a .git by default, so we remove it first
	// and let our scaffold handle it with -b main
	if config.Git {
		// cargo init may have already created a git repo, remove it so we get -b main
		os.RemoveAll(fmt.Sprintf("%s/.git", projectDir))
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		// Remove the git repo cargo created since user didn't ask for it
		os.RemoveAll(fmt.Sprintf("%s/.git", projectDir))
	}

	// Print summary
	fmt.Println()
	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Rust project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	summary.WriteString("  cargo run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
