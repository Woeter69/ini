package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("cobol", &CobolHandler{}) }

type CobolHandler struct{}

func (c *CobolHandler) Name() string { return "COBOL" }
func (c *CobolHandler) Validate() error { return nil }

func (c *CobolHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.cbl...\n", ui.Arrow)

	// Basic free-format COBOL skeleton
	mainCbl := fmt.Sprintf(`       IDENTIFICATION DIVISION.
       PROGRAM-ID. MAIN.

       ENVIRONMENT DIVISION.

       DATA DIVISION.

       PROCEDURE DIVISION.
           DISPLAY "Hello from %s!".
           STOP RUN.
`, config.Name)

	if err := os.WriteFile(filepath.Join(config.Path, "main.cbl"), []byte(mainCbl), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.cbl created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your COBOL project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  cobc -x -free -o main main.cbl\n")
	s.WriteString("  ./main\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
