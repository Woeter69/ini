package handler

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("dart", &DartHandler{}) }

type DartHandler struct{}

func (d *DartHandler) Name() string { return "Dart" }

func (d *DartHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli"}
}

func (d *DartHandler) Validate() error {
	if _, err := exec.LookPath("dart"); err != nil {
		return fmt.Errorf("dart is not installed or not in PATH.\n  Install it: https://dart.dev/get-dart")
	}
	return nil
}

func (d *DartHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	fmt.Printf("  %s Creating Dart project...\n", ui.Arrow)

	// Dart templates: console-full, package, etc.
	template := "console-full"
	if config.Type == "app" {
		template = "package"
	}

	dartCreate := exec.Command("dart", "create", "-t", template, config.Name, "--force")
	parent := filepath.Dir(projectDir)
	if parent == "" {
		parent = "."
	}
	dartCreate.Dir = parent
	if err := dartCreate.Run(); err != nil {
		if err := scaffold.CreateDir(projectDir); err != nil {
			return err
		}
	}
	fmt.Printf("  %s Dart %s project created\n", ui.CheckMark, config.Type)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	d.printSummary(config)
	return nil
}

func (d *DartHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Dart %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  dart run bin/main.dart\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
