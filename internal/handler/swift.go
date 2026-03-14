package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("swift", &SwiftHandler{}) }

type SwiftHandler struct{}

func (s *SwiftHandler) Name() string    { return "Swift" }
func (s *SwiftHandler) Validate() error {
	if _, err := exec.LookPath("swift"); err != nil {
		return fmt.Errorf("swift is not installed or not in PATH.\n  Install it: https://www.swift.org/install/")
	}
	return nil
}

func (sw *SwiftHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Swift package...\n", ui.Arrow)
	swiftInit := exec.Command("swift", "package", "init", "--type", "executable", "--name", config.Name)
	swiftInit.Dir = config.Path
	swiftInit.Stdout = nil
	swiftInit.Stderr = os.Stderr
	if err := swiftInit.Run(); err != nil {
		return fmt.Errorf("failed to init swift package: %w", err)
	}
	fmt.Printf("  %s Swift package initialized\n", ui.CheckMark)

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
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Swift project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  swift run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
