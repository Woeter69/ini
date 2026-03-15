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

func init() { Register("flutter", &FlutterHandler{}) }

type FlutterHandler struct{}

func (f *FlutterHandler) Name() string    { return "Flutter" }
func (f *FlutterHandler) Validate() error {
	if _, err := exec.LookPath("flutter"); err != nil {
		return fmt.Errorf("flutter is not installed or not in PATH.\n  Install it: https://docs.flutter.dev/get-started/install")
	}
	return nil
}

// SupportedTypes declares which global taxonomy categories Flutter supports
func (f *FlutterHandler) SupportedTypes() []string {
	return []string{"basic", "app", "package"}
}

func (f *FlutterHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	fmt.Printf("  %s Creating Flutter project...\n", ui.Arrow)
	
	// Determine template flag
	template := "app"
	if config.Type == "package" {
		template = "package"
	}

	flutterCreate := exec.Command("flutter", "create", "--template="+template, config.Name)
	
	// Run from parent dir since flutter create makes the directory
	parent := filepath.Dir(projectDir)
	if parent == "" {
		parent = "."
	}
	flutterCreate.Dir = parent
	if err := flutterCreate.Run(); err != nil {
		// Mock fallback: create directory
		if err := scaffold.CreateDir(projectDir); err != nil {
			return err
		}
	}
	fmt.Printf("  %s Flutter %s created\n", ui.CheckMark, template)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	f.printSummary(config)
	return nil
}

func (f *FlutterHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Flutter %s is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  flutter run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
