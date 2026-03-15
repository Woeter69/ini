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

func init() {
	Register("go", &GoHandler{})
}

// GoHandler scaffolds Go projects using go mod init.
type GoHandler struct{}

func (g *GoHandler) Name() string {
	return "Go"
}

// SupportedTypes declares which global taxonomy categories Go supports
func (g *GoHandler) SupportedTypes() []string {
	return []string{
		"basic", "app", "cli", "web", "api", "devops", "network", "os", "data", "security",
		"monitor", "stream", "comm", "web3", "lang", "script", "embedded", "math", "stats",
	}
}

func (g *GoHandler) Validate() error {
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go is not installed or not in PATH.\n  Install it: https://go.dev/dl/")
	}
	return nil
}

func (g *GoHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Initialize Go module
	moduleName := config.Name
	fmt.Printf("  %s Initializing Go module (%s)...\n", ui.Arrow, moduleName)
	goMod := exec.Command("go", "mod", "init", moduleName)
	goMod.Dir = projectDir
	if err := goMod.Run(); err != nil {
		return fmt.Errorf("failed to init go module: %w", err)
	}
	fmt.Printf("  %s Go module initialized\n", ui.CheckMark)

	// Determine type path for template
	templatePath := fmt.Sprintf("go/%s/main.go.tmpl", config.Type)
	if config.Type == "" || config.Type == "basic" {
		templatePath = "go/basic/main.go.tmpl"
	}

	// Determine dependencies
	var deps []string
	switch config.Type {
	case "db", "data":
		deps = append(deps, "github.com/mattn/go-sqlite3")
	case "devops":
		deps = append(deps, "github.com/docker/docker/client")
	case "monitor":
		deps = append(deps, "github.com/prometheus/client_golang/prometheus", "github.com/prometheus/client_golang/prometheus/promauto")
	case "web3":
		deps = append(deps, "github.com/ethereum/go-ethereum/ethclient")
	case "api":
		// use web template
	case "cli":
		deps = append(deps, "github.com/spf13/cobra")
	case "math", "stats":
		deps = append(deps, "gonum.org/v1/gonum/mat")
	}

	// 3. Add dependencies if any
	if len(deps) > 0 {
		fmt.Printf("  %s Adding dependencies: %s...\n", ui.Arrow, strings.Join(deps, ", "))
		for _, dep := range deps {
			goGet := exec.Command("go", "get", dep)
			goGet.Dir = projectDir
			if err := goGet.Run(); err != nil {
				return fmt.Errorf("failed to add dependency %s: %w", dep, err)
			}
		}
		fmt.Printf("  %s Dependencies added\n", ui.CheckMark)
	}

	// 4. Create main.go from embedded template
	tmplContent, err := templates.FS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	t, err := template.New("main").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "main.go"), buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}
	fmt.Printf("  %s main.go created\n", ui.CheckMark)

	// 5. Create .gitignore and README
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	// 6. Git
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	g.printSummary(config)
	return nil
}

func (g *GoHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Go %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  go run .\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
