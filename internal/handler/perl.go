package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("perl", &PerlHandler{}) }

type PerlHandler struct{}

func (p *PerlHandler) Name() string    { return "Perl" }
func (p *PerlHandler) Validate() error { return nil }

func (p *PerlHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create main.pl
	fmt.Printf("  %s Creating main.pl...\n", ui.Arrow)
	mainPl := fmt.Sprintf(`#!/usr/bin/env perl
use strict;
use warnings;

# %s

print "Hello from %s!\n";
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.pl"), []byte(mainPl), 0o755); err != nil {
		return err
	}
	fmt.Printf("  %s main.pl created\n", ui.CheckMark)

	// Create cpanfile
	fmt.Printf("  %s Creating cpanfile...\n", ui.Arrow)
	cpanfile := "# Add your dependencies here\n# requires 'Mojolicious';\n# requires 'DBI';\n"
	if err := os.WriteFile(filepath.Join(config.Path, "cpanfile"), []byte(cpanfile), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s cpanfile created\n", ui.CheckMark)

	// Create lib/ directory
	os.MkdirAll(filepath.Join(config.Path, "lib"), 0o755)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Perl project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  perl main.pl\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
