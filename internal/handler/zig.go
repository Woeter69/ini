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
	Register("zig", &ZigHandler{})
}

// ZigHandler scaffolds Zig projects using zig init.
type ZigHandler struct{}

func (z *ZigHandler) Name() string {
	return "Zig"
}

// SupportedTypes declares which global taxonomy categories Zig supports
func (z *ZigHandler) SupportedTypes() []string {
	return []string{"basic", "cli", "embedded", "web", "game"}
}

func (z *ZigHandler) Validate() error {
	_, err := exec.LookPath("zig")
	if err != nil {
		return fmt.Errorf("zig is not installed or not in PATH.\n  Install it: https://ziglang.org/download/")
	}
	return nil
}

func (z *ZigHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Run zig init (provides basic structure)
	fmt.Printf("  %s Initializing Zig project...\n", ui.Arrow)
	zigInit := exec.Command("zig", "init")
	zigInit.Dir = projectDir
	if err := zigInit.Run(); err != nil {
		return fmt.Errorf("failed to init zig project: %w", err)
	}
	fmt.Printf("  %s Zig project structure initialized\n", ui.CheckMark)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("zig/%s/main.zig.tmpl", typeDir)
	buildTmplPath := fmt.Sprintf("zig/%s/build.zig.tmpl", typeDir)

	// 3. Overwrite src/main.zig
	fmt.Printf("  %s Generating src/main.zig...\n", ui.Arrow)
	if err := z.processTemplate(config, mainTmplPath, filepath.Join(projectDir, "src", "main.zig")); err != nil {
		return err
	}

	// 4. Overwrite build.zig (fallback to basic if not found)
	fmt.Printf("  %s Configuring build.zig...\n", ui.Arrow)
	if err := z.processTemplate(config, buildTmplPath, filepath.Join(projectDir, "build.zig")); err != nil {
		// If specific build.zig doesn't exist, we skip overwriting (keep default from zig init)
	}

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
	}

	z.printSummary(config)
	return nil
}

func (z *ZigHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback for build.zig
		if strings.HasSuffix(tmplPath, "build.zig.tmpl") {
			tmplPath = "zig/basic/build.zig.tmpl"
			content, err = templates.FS.ReadFile(tmplPath)
		}
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
		}
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

func (z *ZigHandler) printSummary(config ProjectConfig) {
	fmt.Println()
	relPath, _ := filepath.Rel(".", config.Path)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Zig %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  zig build run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()
}
