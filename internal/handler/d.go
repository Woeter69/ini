package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("d", &DHandler{}) }

type DHandler struct{}

func (d *DHandler) Name() string { return "D" }
func (d *DHandler) Validate() error {
	if _, err := exec.LookPath("dub"); err != nil {
		return fmt.Errorf("dub is not installed or not in PATH.\n  Install it: https://dlang.org/download.html")
	}
	return nil
}

func (d *DHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing D project with dub...\n", ui.Arrow)
	dubInit := exec.Command("dub", "init", "-n") // -n non-interactive
	dubInit.Dir = config.Path
	dubInit.Stdout = nil
	dubInit.Stderr = os.Stderr
	if err := dubInit.Run(); err != nil {
		return fmt.Errorf("failed to init dub project: %w", err)
	}
	fmt.Printf("  %s dub project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your D project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  dub run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
