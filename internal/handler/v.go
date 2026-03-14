package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("v", &VHandler{}) }

type VHandler struct{}

func (v *VHandler) Name() string { return "V" }
func (v *VHandler) Validate() error {
	if _, err := exec.LookPath("v"); err != nil {
		return fmt.Errorf("v compiler is not installed or not in PATH.\n  Install it: https://vlang.io/")
	}
	return nil
}

func (v *VHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing V project...\n", ui.Arrow)
	vInit := exec.Command("v", "init")
	vInit.Dir = config.Path
	vInit.Stdout = nil
	vInit.Stderr = os.Stderr
	if err := vInit.Run(); err != nil {
		return fmt.Errorf("failed to init v project: %w", err)
	}
	fmt.Printf("  %s v project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git {
		os.RemoveAll(config.Path + "/.git") // v init creates git
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your V project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  v run .\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
