package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("julia", &JuliaHandler{}) }

type JuliaHandler struct{}

func (j *JuliaHandler) Name() string    { return "Julia" }
func (j *JuliaHandler) Validate() error { return nil }

func (j *JuliaHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.jl...\n", ui.Arrow)
	mainJl := fmt.Sprintf("# %s\n\nprintln(\"Hello from %s!\")\n", config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.jl"), []byte(mainJl), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.jl created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Julia project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  julia main.jl\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
