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

func init() { Register("bun", &BunHandler{}) }

type BunHandler struct{}

func (b *BunHandler) Name() string { return "JavaScript/TypeScript (Bun)" }

func (b *BunHandler) SupportedTypes() []string {
	return []string{"basic", "app", "web", "api", "cli", "os", "network", "data"}
}

func (b *BunHandler) Validate() error {
	if _, err := exec.LookPath("bun"); err != nil {
		return fmt.Errorf("bun is not installed or not in PATH.\n  Install it: curl -fsSL https://bun.sh/install | bash")
	}
	return nil
}

func (b *BunHandler) Init(config ProjectConfig) error {
	projectDir := config.Path
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Handle Web Frameworks specially
	if config.Type == "web" && config.Framework != "" {
		return b.scaffoldWeb(config)
	}

	// Handle other custom templates
	if config.Type != "basic" && config.Type != "" {
		typeDir := config.Type
		if typeDir == "cli" { typeDir = "script" }
		if typeDir == "data" { typeDir = "db" }
		if typeDir == "app" || typeDir == "api" { typeDir = "web" } // fallback

		templateDir := filepath.Join("bun", typeDir)
		if _, err := templates.FS.ReadDir(templateDir); err == nil {
			return b.scaffoldTemplate(config, templateDir)
		}
	}

	// Default Bun Init
	fmt.Printf("  %s Initializing Bun project...\n", ui.Arrow)
	bunInit := exec.Command("bun", "init", "-y")
	bunInit.Dir = projectDir
	bunInit.Stdout = nil
	bunInit.Stderr = os.Stderr
	if err := bunInit.Run(); err != nil {
		return fmt.Errorf("failed to init bun project: %w", err)
	}
	fmt.Printf("  %s bun project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return err
	}

	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	b.printSummary(config)
	return nil
}

func (b *BunHandler) scaffoldWeb(config ProjectConfig) error {
	fmt.Printf("  %s Scaffolding %s (%s) project...\n", ui.Arrow, config.Framework, config.Variant)

	switch config.Framework {
	case "next":
		nextArgs := []string{"x", "create-next-app@latest", ".", "--tailwind", "--eslint", "--app", "--src-dir", "--use-bun", "--no-git"}
		if config.Variant == "ts" {
			nextArgs = append(nextArgs, "--ts")
		} else {
			nextArgs = append(nextArgs, "--js")
		}
		cmd := exec.Command("bun", nextArgs...)
		cmd.Dir = config.Path
		cmd.Stdin = strings.NewReader("y\n")
		return cmd.Run()
	case "react":
		tDir := "bun/react-jsx"
		if config.Variant == "ts" {
			tDir = "bun/react-tsx"
		}
		if err := b.scaffoldTemplate(config, tDir); err != nil {
			return err
		}
		fmt.Printf("  %s Installing dependencies...\n", ui.Arrow)
		return exec.Command("bun", "install").Run()
	case "vue":
		tDir := "bun/vue-js"
		if config.Variant == "ts" {
			tDir = "bun/vue-ts"
		}
		if err := b.scaffoldTemplate(config, tDir); err != nil {
			return err
		}
		fmt.Printf("  %s Installing dependencies...\n", ui.Arrow)
		return exec.Command("bun", "install").Run()
	case "svelte":
		tDir := "bun/svelte-js"
		if config.Variant == "ts" {
			tDir = "bun/svelte-ts"
		}
		if err := b.scaffoldTemplate(config, tDir); err != nil {
			return err
		}
		fmt.Printf("  %s Installing dependencies...\n", ui.Arrow)
		if err := exec.Command("bun", "install").Run(); err != nil {
			return err
		}
	case "solid":
		tDir := "bun/solid-jsx"
		if config.Variant == "ts" {
			tDir = "bun/solid-tsx"
		}
		if err := b.scaffoldTemplate(config, tDir); err != nil {
			return err
		}
		fmt.Printf("  %s Installing dependencies...\n", ui.Arrow)
		if err := exec.Command("bun", "install").Run(); err != nil {
			return err
		}
	case "angular":
		// Angular CLI is TS by default
		fmt.Printf("  %s Using external Angular CLI (Angular requires many auxiliary tools)...\n", ui.Arrow)
		cmd := exec.Command("bun", "x", "@angular/cli", "new", config.Name, "--directory", ".", "--skip-git")
		cmd.Dir = config.Path
		if err := cmd.Run(); err != nil {
			return err
		}
	case "express":
		if err := b.scaffoldTemplate(config, "bun/express"); err != nil {
			return err
		}
		exec.Command("bun", "init", "-y").Run()
		fmt.Printf("  %s Adding express...\n", ui.Arrow)
		cmd := exec.Command("bun", "add", "express")
		cmd.Dir = config.Path
		if err := cmd.Run(); err != nil {
			return err
		}
	case "vanilla":
		tDir := "bun/vanilla-js"
		if config.Variant == "ts" {
			tDir = "bun/vanilla-ts"
		}
		if err := b.scaffoldTemplate(config, tDir); err != nil {
			return err
		}
	}

	b.printSummary(config)
	return nil
}

func (b *BunHandler) scaffoldTemplate(config ProjectConfig, templateDir string) error {
	entries, err := templates.FS.ReadDir(templateDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		subPath := filepath.Join(templateDir, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(filepath.Join(config.Path, entry.Name()), 0755); err != nil {
				return err
			}
			if err := b.copyRecursive(config, subPath, entry.Name()); err != nil {
				return err
			}
			continue
		}

		if err := b.processFile(config, subPath, entry.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (b *BunHandler) copyRecursive(config ProjectConfig, srcDir, destPrefix string) error {
	entries, err := templates.FS.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destPrefix, entry.Name())
		if entry.IsDir() {
			if err := os.MkdirAll(filepath.Join(config.Path, destPath), 0755); err != nil {
				return err
			}
			if err := b.copyRecursive(config, srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := b.processFile(config, srcPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BunHandler) printSummary(config ProjectConfig) error {
	fmt.Println()
	s := strings.Builder{}
	title := "JavaScript/TypeScript"
	if config.Framework != "" {
		title = config.Framework
	}
	s.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your %s project is ready!", title)))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	
	switch config.Framework {
	case "next", "react-js", "react-ts", "vue":
		s.WriteString("  bun run dev\n")
	case "express":
		s.WriteString("  bun run src/index.js\n")
	case "vanilla-js", "vanilla-ts":
		s.WriteString("  # Open index.html in your browser\n")
	default:
		s.WriteString("  bun run index.ts\n")
	}
	
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}

func (b *BunHandler) processFile(config ProjectConfig, srcPath, destRelPath string) error {
	// Skip files based on variant if they are code files
	// e.g. if variant is "js", skip ".ts" files, and vice-versa.
	// We only apply this to the final destination extension (after stripping .tmpl)
	targetName := strings.TrimSuffix(destRelPath, ".tmpl")
	ext := filepath.Ext(targetName)

	if config.Variant != "" {
		if config.Variant == "js" && (ext == ".ts" || ext == ".tsx") {
			return nil
		}
		if config.Variant == "ts" && (ext == ".js" || ext == ".jsx") {
			// Skip .js/.jsx in TS projects, UNLESS it's a specific requirement (next.config.js etc.)
			// For our utility templates (index.js), we certainly want to skip.
			return nil
		}
	}

	content, err := templates.FS.ReadFile(srcPath)
	if err != nil {
		return err
	}

	// Execute template if it has .tmpl extension (either in source or dest)
	if strings.HasSuffix(srcPath, ".tmpl") || strings.HasSuffix(destRelPath, ".tmpl") {
		t, err := template.New(filepath.Base(srcPath)).Delims("[[", "]]").Parse(string(content))
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, config); err != nil {
			return err
		}
		content = buf.Bytes()
	}

	return os.WriteFile(filepath.Join(config.Path, targetName), content, 0644)
}

