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

func init() { Register("php", &PHPHandler{}) }

type PHPHandler struct{}

func (p *PHPHandler) Name() string    { return "PHP" }
func (p *PHPHandler) Validate() error {
	if _, err := exec.LookPath("php"); err != nil {
		return fmt.Errorf("php is not installed or not in PATH.\n  Install it: https://www.php.net/downloads")
	}
	return nil
}

func (p *PHPHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create src/ directory
	os.MkdirAll(filepath.Join(config.Path, "src"), 0o755)

	// Create main.php
	fmt.Printf("  %s Creating main.php...\n", ui.Arrow)
	mainPhp := fmt.Sprintf(`<?php

declare(strict_types=1);

// %s

echo "Hello from %s!" . PHP_EOL;
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.php"), []byte(mainPhp), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.php created\n", ui.CheckMark)

	// Create composer.json
	fmt.Printf("  %s Creating composer.json...\n", ui.Arrow)
	composerJson := fmt.Sprintf(`{
    "name": "%s/%s",
    "description": "A PHP project",
    "type": "project",
    "require": {
        "php": ">=8.2"
    },
    "autoload": {
        "psr-4": {
            "App\\": "src/"
        }
    }
}
`, config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "composer.json"), []byte(composerJson), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s composer.json created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your PHP project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  php main.php\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
