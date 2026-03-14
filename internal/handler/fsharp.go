package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("fsharp", &FSharpHandler{}) }

type FSharpHandler struct{}

func (f *FSharpHandler) Name() string { return "F#" }
func (f *FSharpHandler) Validate() error {
	if _, err := exec.LookPath("dotnet"); err != nil {
		return fmt.Errorf("dotnet is not installed or not in PATH.\n  Install it: https://dotnet.microsoft.com/download")
	}
	return nil
}

func (f *FSharpHandler) Init(config ProjectConfig) error {
	fmt.Printf("  %s Creating F# project with dotnet...\n", ui.Arrow)
	dotnetNew := exec.Command("dotnet", "new", "console", "-lang", "F#", "-n", config.Name, "-o", config.Path, "--force")
	dotnetNew.Stdout = nil
	dotnetNew.Stderr = os.Stderr
	if err := dotnetNew.Run(); err != nil {
		return fmt.Errorf("failed to create dotnet F# project: %w", err)
	}
	fmt.Printf("  %s F# project created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your F# project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  dotnet run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
