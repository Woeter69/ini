package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("pascal", &PascalHandler{}) }

type PascalHandler struct{}

func (p *PascalHandler) Name() string { return "Pascal" }
func (p *PascalHandler) Validate() error { return nil }

func (p *PascalHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.pas...\n", ui.Arrow)

	mainPas := fmt.Sprintf(`program %s;

begin
  writeln('Hello from %s!');
end.
`, config.Name, config.Name)
	// Replace non-alphanumeric chars for program name (Pascal gets strict about program identifiers)
	safeProgramName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, config.Name)
	
	mainPas = fmt.Sprintf(`program %s;

begin
  writeln('Hello from %s!');
end.
`, safeProgramName, config.Name)

	if err := os.WriteFile(filepath.Join(config.Path, "main.pas"), []byte(mainPas), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.pas created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Pascal project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  fpc main.pas\n")
	s.WriteString("  ./main\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
