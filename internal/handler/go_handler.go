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
	Register("go", &GoHandler{})
}

// GoHandler scaffolds Go projects using go mod init.
type GoHandler struct{}

func (g *GoHandler) Name() string {
	return "Go"
}

func (g *GoHandler) Validate() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go is not installed or not in PATH.\n  Install it: https://go.dev/dl/")
	}
	return nil
}

func (g *GoHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Initialize Go module
	moduleName := config.Name
	// If it looks like a plain name, prefix with a sensible default
	if !strings.Contains(moduleName, "/") {
		moduleName = config.Name
	}
	fmt.Printf("  %s Initializing Go module (%s)...\n", ui.Arrow, moduleName)
	goMod := exec.Command("go", "mod", "init", moduleName)
	goMod.Dir = projectDir
	goMod.Stdout = nil
	goMod.Stderr = os.Stderr
	if err := goMod.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}
	fmt.Printf("  %s Go module initialized\n", ui.CheckMark)

	// 3. Create main.go
	fmt.Printf("  %s Creating main.go...\n", ui.Arrow)
	mainGo := fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, config.Name)
	if err := os.WriteFile(filepath.Join(projectDir, "main.go"), []byte(mainGo), 0o644); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	fmt.Printf("  %s main.go created\n", ui.CheckMark)

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
	summary.WriteString(ui.SuccessStyle.Render("🚀 Your Go project is ready!"))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  go run .\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
