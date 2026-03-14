package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("crystal", &CrystalHandler{}) }

type CrystalHandler struct{}

func (c *CrystalHandler) Name() string { return "Crystal" }
func (c *CrystalHandler) Validate() error {
	if _, err := exec.LookPath("crystal"); err != nil {
		return fmt.Errorf("crystal is not installed or not in PATH.\n  Install it: https://crystal-lang.org/install/")
	}
	return nil
}

func (c *CrystalHandler) Init(config ProjectConfig) error {
	fmt.Printf("  %s Initializing Crystal project...\n", ui.Arrow)
	parent := config.Path[:len(config.Path)-len(config.Name)-1]
	if parent == "" { parent = "." }
	crystalInit := exec.Command("crystal", "init", "app", config.Name)
	crystalInit.Dir = parent
	crystalInit.Stdout = nil
	crystalInit.Stderr = os.Stderr
	if err := crystalInit.Run(); err != nil {
		return fmt.Errorf("failed to init crystal project: %w", err)
	}
	fmt.Printf("  %s crystal project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git {
		os.RemoveAll(config.Path + "/.git") // crystal init app creates git
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Crystal project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString(fmt.Sprintf("  crystal run src/%s.cr\n", config.Name))
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
