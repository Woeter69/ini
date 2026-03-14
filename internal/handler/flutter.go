package handler

import (
	"fmt"
	"os"
	"os/exec"
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

func (f *FlutterHandler) Init(config ProjectConfig) error {
	fmt.Printf("  %s Creating Flutter project...\n", ui.Arrow)
	flutterCreate := exec.Command("flutter", "create", config.Name)
	// Run from parent dir since flutter create makes the directory
	parent := config.Path[:len(config.Path)-len(config.Name)-1]
	if parent == "" { parent = "." }
	flutterCreate.Dir = parent
	flutterCreate.Stdout = nil
	flutterCreate.Stderr = os.Stderr
	if err := flutterCreate.Run(); err != nil {
		return fmt.Errorf("failed to create flutter project: %w", err)
	}
	fmt.Printf("  %s Flutter project created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil {
		return err
	}
	if config.Git {
		os.RemoveAll(config.Path + "/.git")
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Flutter project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  flutter run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
