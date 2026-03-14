package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("elixir", &ElixirHandler{}) }

type ElixirHandler struct{}

func (e *ElixirHandler) Name() string { return "Elixir" }
func (e *ElixirHandler) Validate() error {
	if _, err := exec.LookPath("mix"); err != nil {
		return fmt.Errorf("mix (Elixir) is not installed or not in PATH.\n  Install it: https://elixir-lang.org/install.html")
	}
	return nil
}

func (e *ElixirHandler) Init(config ProjectConfig) error {
	// mix new creates the directory itself
	fmt.Printf("  %s Creating Elixir project with mix...\n", ui.Arrow)
	parent := config.Path[:len(config.Path)-len(config.Name)-1]
	if parent == "" { parent = "." }
	mixNew := exec.Command("mix", "new", config.Name)
	mixNew.Dir = parent
	mixNew.Stdout = nil
	mixNew.Stderr = os.Stderr
	if err := mixNew.Run(); err != nil {
		return fmt.Errorf("failed to create mix project: %w", err)
	}
	fmt.Printf("  %s Mix project created\n", ui.CheckMark)

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
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Elixir project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  iex -S mix\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
