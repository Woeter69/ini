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

func init() { Register("ruby", &RubyHandler{}) }

type RubyHandler struct{}

func (r *RubyHandler) Name() string    { return "Ruby" }
func (r *RubyHandler) Validate() error {
	if _, err := exec.LookPath("ruby"); err != nil {
		return fmt.Errorf("ruby is not installed or not in PATH.\n  Install it: https://www.ruby-lang.org/en/downloads/")
	}
	return nil
}

func (r *RubyHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create main.rb
	fmt.Printf("  %s Creating main.rb...\n", ui.Arrow)
	mainRb := fmt.Sprintf("# %s\n\nputs \"Hello from %s!\"\n", config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.rb"), []byte(mainRb), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.rb created\n", ui.CheckMark)

	// Create Gemfile
	fmt.Printf("  %s Creating Gemfile...\n", ui.Arrow)
	gemfile := "# frozen_string_literal: true\n\nsource \"https://rubygems.org\"\n\n# Add your gems here\n# gem \"httparty\"\n"
	if err := os.WriteFile(filepath.Join(config.Path, "Gemfile"), []byte(gemfile), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s Gemfile created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Ruby project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  ruby main.rb\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
