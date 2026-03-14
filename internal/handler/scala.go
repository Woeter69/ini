package handler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("scala", &ScalaHandler{}) }

type ScalaHandler struct{}

func (s *ScalaHandler) Name() string    { return "Scala" }
func (s *ScalaHandler) Validate() error {
	if _, err := exec.LookPath("scala-cli"); err != nil {
		return fmt.Errorf("scala-cli is not installed or not in PATH.\n  Install it: https://scala-cli.virtuslab.org/install")
	}
	return nil
}

func (sc *ScalaHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create main.scala
	fmt.Printf("  %s Creating main.scala...\n", ui.Arrow)
	mainScala := fmt.Sprintf(`//> using scala 3
//> using option -Werror

@main def run(): Unit =
  println("Hello from %s!")
`, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.scala"), []byte(mainScala), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.scala created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Scala project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  scala-cli run .\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
