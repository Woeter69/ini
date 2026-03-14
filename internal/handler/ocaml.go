package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("ocaml", &OCamlHandler{}) }

type OCamlHandler struct{}

func (o *OCamlHandler) Name() string { return "OCaml" }
func (o *OCamlHandler) Validate() error {
	if _, err := exec.LookPath("dune"); err != nil {
		return fmt.Errorf("dune is not installed or not in PATH.\n  Install it: opam install dune")
	}
	return nil
}

func (o *OCamlHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing OCaml project with dune...\n", ui.Arrow)
	duneInit := exec.Command("dune", "init", "project", config.Name)
	parent := config.Path[:len(config.Path)-len(config.Name)-1]
	if parent == "" { parent = "." }
	duneInit.Dir = parent
	duneInit.Stdout = nil
	duneInit.Stderr = os.Stderr
	if err := duneInit.Run(); err != nil {
		return fmt.Errorf("failed to init dune project: %w", err)
	}
	fmt.Printf("  %s dune project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git {
		os.RemoveAll(config.Path + "/.git")
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your OCaml project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  dune build\n")
	s.WriteString("  dune exec ./bin/main.exe\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
