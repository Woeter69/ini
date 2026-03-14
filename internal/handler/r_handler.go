package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("r", &RHandler{}) }

type RHandler struct{}

func (r *RHandler) Name() string    { return "R" }
func (r *RHandler) Validate() error { return nil }

func (r *RHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create main.R
	fmt.Printf("  %s Creating main.R...\n", ui.Arrow)
	mainR := fmt.Sprintf("# %s\n\ncat(\"Hello from %s!\\n\")\n", config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.R"), []byte(mainR), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.R created\n", ui.CheckMark)

	// Create R/ directory for functions
	os.MkdirAll(filepath.Join(config.Path, "R"), 0o755)

	// Create .Rprofile
	fmt.Printf("  %s Creating .Rprofile...\n", ui.Arrow)
	rprofile := "# Project-level R settings\n# options(repos = c(CRAN = \"https://cloud.r-project.org\"))\n"
	if err := os.WriteFile(filepath.Join(config.Path, ".Rprofile"), []byte(rprofile), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s .Rprofile created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your R project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  Rscript main.R\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
