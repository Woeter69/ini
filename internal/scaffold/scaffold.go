package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Woeter69/ini/internal/templates"
	"github.com/Woeter69/ini/internal/ui"
)

// CreateDir creates the project directory.
func CreateDir(projectDir string) error {
	fmt.Printf("  %s Creating project directory...\n", ui.Arrow)
	return os.MkdirAll(projectDir, 0o755)
}

// WriteGitignore writes a language-specific .gitignore.
func WriteGitignore(projectDir, language string) error {
	content, ok := templates.Gitignore[language]
	if !ok {
		return fmt.Errorf("no .gitignore template for %q", language)
	}
	fmt.Printf("  %s Creating .gitignore...\n", ui.Arrow)
	if err := os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s .gitignore created\n", ui.CheckMark)
	return nil
}

// WriteReadme writes a language-specific README.md.
func WriteReadme(projectDir, projectName, language string) error {
	content := templates.Readme(projectName, language)
	fmt.Printf("  %s Creating README.md...\n", ui.Arrow)
	if err := os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s README.md created\n", ui.CheckMark)
	return nil
}

// InitGit initializes a git repo with -b main if git is available.
func InitGit(projectDir string) error {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Printf("  %s git not found, skipping\n", ui.WarningStyle.Render("⚠"))
		return nil
	}

	gitDir := filepath.Join(projectDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		// Already a git repo
		return nil
	}

	fmt.Printf("  %s Initializing git repository (branch: main)...\n", ui.Arrow)
	gitInit := exec.Command("git", "init", "-b", "main")
	gitInit.Dir = projectDir
	gitInit.Stdout = nil
	gitInit.Stderr = nil
	if err := gitInit.Run(); err != nil {
		return fmt.Errorf("failed to init git: %w", err)
	}
	fmt.Printf("  %s Git repository initialized (main)\n", ui.CheckMark)
	return nil
}
