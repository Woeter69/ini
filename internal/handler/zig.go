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
	Register("zig", &ZigHandler{})
}

// ZigHandler scaffolds Zig projects using zig init.
type ZigHandler struct{}

func (z *ZigHandler) Name() string {
	return "Zig"
}

func (z *ZigHandler) Validate() error {
	_, err := exec.LookPath("zig")
	if err != nil {
		return fmt.Errorf("zig is not installed or not in PATH.\n  Install it: https://ziglang.org/download/")
	}
	return nil
}

func (z *ZigHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Run zig init
	fmt.Printf("  %s Initializing Zig project...\n", ui.Arrow)
	zigInit := exec.Command("zig", "init")
	zigInit.Dir = projectDir
	zigInit.Stdout = nil
	zigInit.Stderr = os.Stderr
	if err := zigInit.Run(); err != nil {
		return fmt.Errorf("failed to init zig project: %w", err)
	}
	fmt.Printf("  %s Zig project initialized\n", ui.CheckMark)

	// 3. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 4. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 5. Initialize git repo (if --git flag is set)
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	// Print summary
	fmt.Println()
	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Zig project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	summary.WriteString("  zig build run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
