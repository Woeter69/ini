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

func init() { Register("kotlin", &KotlinHandler{}) }

type KotlinHandler struct{}

func (k *KotlinHandler) Name() string { return "Kotlin" }

// SupportedTypes declares which global taxonomy categories Kotlin supports
func (k *KotlinHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "api", "cli", "data", "desktop", "ai"}
}

func (k *KotlinHandler) Validate() error {
	if _, err := exec.LookPath("gradle"); err != nil {
		return fmt.Errorf("gradle is not installed or not in PATH.\n  Install it: https://gradle.org/install/")
	}
	return nil
}

func (k *KotlinHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Kotlin project with Gradle...\n", ui.Arrow)
	// Non-interactive gradle init
	gradleInit := exec.Command("gradle", "init",
		"--type", "kotlin-application",
		"--dsl", "kotlin",
		"--test-framework", "kotlintest",
		"--project-name", config.Name,
		"--package", config.Name,
		"--no-split-project",
		"--java-version", "21")
	gradleInit.Dir = projectDir
	if err := gradleInit.Run(); err != nil {
		return fmt.Errorf("failed to init gradle project: %w", err)
	}
	fmt.Printf("  %s Gradle Kotlin project initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" || typeDir == "cli" {
		typeDir = "basic"
	}
	if typeDir == "api" {
		typeDir = "web"
	}
	if typeDir == "db" {
		typeDir = "data"
	}
	appTmplPath := fmt.Sprintf("kotlin/%s/App.kt.tmpl", typeDir)
	buildTmplPath := fmt.Sprintf("kotlin/%s/build.gradle.kts.tmpl", typeDir)

	// In Gradle "kotlin-application" with no-split-project,
	// files are in app/src/main/kotlin/PACKAGE/App.kt
	appPath := filepath.Join(projectDir, "app", "src", "main", "kotlin", config.Name, "App.kt")
	buildPath := filepath.Join(projectDir, "app", "build.gradle.kts")

	// 2. Overwrite build.gradle.kts
	fmt.Printf("  %s Configuring build.gradle.kts...\n", ui.Arrow)
	if err := k.processTemplate(config, buildTmplPath, buildPath); err != nil {
		return err
	}

	// 3. Overwrite App.kt
	fmt.Printf("  %s Generating App.kt...\n", ui.Arrow)
	if err := k.processTemplate(config, appTmplPath, appPath); err != nil {
		return err
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
	}

	k.printSummary(config)
	return nil
}

func (k *KotlinHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
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

func (k *KotlinHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Kotlin %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  gradle run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
