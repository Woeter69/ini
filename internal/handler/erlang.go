package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("erlang", &ErlangHandler{}) }

type ErlangHandler struct{}

func (e *ErlangHandler) Name() string { return "Erlang" }
func (e *ErlangHandler) Validate() error {
	if _, err := exec.LookPath("rebar3"); err != nil {
		return fmt.Errorf("rebar3 is not installed or not in PATH.\n  Install it: https://rebar3.org/docs/getting-started/")
	}
	return nil
}

func (e *ErlangHandler) Init(config ProjectConfig) error {
	// rebar3 new app creates the directory itself
	fmt.Printf("  %s Creating Erlang project with rebar3...\n", ui.Arrow)
	parent := config.Path[:len(config.Path)-len(config.Name)-1]
	if parent == "" { parent = "." }
	rebar := exec.Command("rebar3", "new", "app", config.Name)
	rebar.Dir = parent
	rebar.Stdout = nil
	rebar.Stderr = os.Stderr
	if err := rebar.Run(); err != nil {
		return fmt.Errorf("failed to create rebar3 project: %w", err)
	}
	fmt.Printf("  %s rebar3 project created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Erlang project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  rebar3 shell\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
