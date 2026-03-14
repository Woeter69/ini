package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("ada", &AdaHandler{}) }

type AdaHandler struct{}

func (a *AdaHandler) Name() string { return "Ada" }
func (a *AdaHandler) Validate() error { return nil }

func (a *AdaHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.adb...\n", ui.Arrow)

	// Ada program names must be valid identifiers
	safeProgramName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, config.Name)
	
	mainAdb := fmt.Sprintf(`with Ada.Text_IO; use Ada.Text_IO;

procedure %s is
begin
   Put_Line ("Hello from %s!");
end %s;
`, safeProgramName, config.Name, safeProgramName)

	if err := os.WriteFile(filepath.Join(config.Path, "main.adb"), []byte(mainAdb), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.adb created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Ada project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  gnatmake main.adb\n")
	s.WriteString("  ./main\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
