package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("nim", &NimHandler{}) }

type NimHandler struct{}

func (n *NimHandler) Name() string { return "Nim" }
func (n *NimHandler) Validate() error {
	if _, err := exec.LookPath("nimble"); err != nil {
		return fmt.Errorf("nimble is not installed or not in PATH.\n  Install it: https://nim-lang.org/install.html")
	}
	return nil
}

func (n *NimHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Nim project with nimble...\n", ui.Arrow)
	nimbleInit := exec.Command("nimble", "init", "-y") // -y auto-answers yes
	nimbleInit.Dir = config.Path
	nimbleInit.Stdout = nil
	nimbleInit.Stderr = os.Stderr
	if err := nimbleInit.Run(); err != nil {
		return fmt.Errorf("failed to init nimble project: %w", err)
	}
	fmt.Printf("  %s nimble project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Nim project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	// nimble init typically creates src/{name}.nim
	s.WriteString(fmt.Sprintf("  nim c -r src/%s.nim\n", config.Name))
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
