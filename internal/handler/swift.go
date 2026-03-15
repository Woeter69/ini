package handler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/templates"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("swift", &SwiftHandler{}) }

type SwiftHandler struct{}

func (s *SwiftHandler) Name() string { return "Swift" }

// SupportedTypes declares which global taxonomy categories Swift supports
func (s *SwiftHandler) SupportedTypes() []string {
	return []string{"basic", "cli", "server", "ios"}
}

func (s *SwiftHandler) Validate() error {
	if _, err := exec.LookPath("swift"); err != nil {
		return fmt.Errorf("swift is not installed or not in PATH.\n  Install it: https://www.swift.org/install/")
	}
	return nil
}

func (sw *SwiftHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Swift package...\n", ui.Arrow)
	initType := "executable"
	if config.Type == "ios" {
		initType = "library"
	}
	swiftInit := exec.Command("swift", "package", "init", "--type", initType, "--name", config.Name)
	swiftInit.Dir = projectDir
	if err := swiftInit.Run(); err != nil {
		return fmt.Errorf("failed to init swift package: %w", err)
	}
	fmt.Printf("  %s Swift package structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	packageTmplPath := fmt.Sprintf("swift/%s/Package.swift.tmpl", typeDir)
	mainTmplPath := fmt.Sprintf("swift/%s/main.swift.tmpl", typeDir)

	// 1. Overwrite Package.swift
	fmt.Printf("  %s Configuring Package.swift...\n", ui.Arrow)
	if err := sw.processTemplate(config, packageTmplPath, filepath.Join(projectDir, "Package.swift")); err != nil {
		return err
	}

	// 2. Overwrite Sources/main.swift (if executable)
	if config.Type != "ios" {
		fmt.Printf("  %s Generating main.swift...\n", ui.Arrow)
		// swift package init creates Sources/main.swift or Sources/<Name>.swift
		// We'll try to find any .swift file in Sources and replace it
		files, _ := filepath.Glob(filepath.Join(projectDir, "Sources", "*.swift"))
		for _, f := range files {
			if err := sw.processTemplate(config, mainTmplPath, f); err != nil {
				return err
			}
		}
	}

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	sw.printSummary(config)
	return nil
}

func (sw *SwiftHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	t, err := template.New(filepath.Base(tmplPath)).Delims("[[", "]]").Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(destPath, buf.Bytes(), 0o644)
}

func (sw *SwiftHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Swift %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  swift run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
