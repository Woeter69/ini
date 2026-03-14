package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("groovy", &GroovyHandler{}) }

type GroovyHandler struct{}

func (g *GroovyHandler) Name() string { return "Groovy" }
func (g *GroovyHandler) Validate() error { return nil }

func (g *GroovyHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.groovy...\n", ui.Arrow)

	mainGroovy := fmt.Sprintf(`// %s

println "Hello from %s!"
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.groovy"), []byte(mainGroovy), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.groovy created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Groovy project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  groovy main.groovy\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
