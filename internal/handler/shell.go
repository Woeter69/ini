package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
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

	// 2. Create main script
	scriptName := config.Name + ".sh"
	fmt.Printf("  %s Creating %s...\n", ui.Arrow, scriptName)
	mainSh := fmt.Sprintf(`#!/usr/bin/env bash
set -euo pipefail

# %s — A shell script project.

main() {
    echo "Hello from %s!"
}

main "$@"
`, config.Name, config.Name)
	scriptPath := filepath.Join(projectDir, scriptName)
	if err := os.WriteFile(scriptPath, []byte(mainSh), 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", scriptName, err)
	}
	fmt.Printf("  %s %s created (executable)\n", ui.CheckMark, scriptName)

	// 3. Create lib/utils.sh
	fmt.Printf("  %s Creating lib/utils.sh...\n", ui.Arrow)
	utilsSh := `#!/usr/bin/env bash
# Utility functions — source this from your main script:
#   source "$(dirname "$0")/lib/utils.sh"

log_info() {
    echo "[INFO] $*"
}

log_error() {
    echo "[ERROR] $*" >&2
}

die() {
    log_error "$@"
    exit 1
}
`
	if err := os.WriteFile(filepath.Join(projectDir, "lib", "utils.sh"), []byte(utilsSh), 0o755); err != nil {
		return fmt.Errorf("failed to create lib/utils.sh: %w", err)
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
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Shell project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString(fmt.Sprintf("  ./%s.sh\n", config.Name))

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
