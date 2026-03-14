package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("haskell", &HaskellHandler{}) }

type HaskellHandler struct{}

func (h *HaskellHandler) Name() string { return "Haskell" }
func (h *HaskellHandler) Validate() error {
	if _, err := exec.LookPath("cabal"); err != nil {
		return fmt.Errorf("cabal is not installed or not in PATH.\n  Install it: https://www.haskell.org/ghcup/")
	}
	return nil
}

func (hs *HaskellHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Haskell project with cabal...\n", ui.Arrow)
	cabalInit := exec.Command("cabal", "init", "--non-interactive", "--package-name", config.Name)
	cabalInit.Dir = config.Path
	cabalInit.Stdout = nil
	cabalInit.Stderr = os.Stderr
	if err := cabalInit.Run(); err != nil {
		return fmt.Errorf("failed to init cabal project: %w", err)
	}
	fmt.Printf("  %s Cabal project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Haskell project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  cabal run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
