package handler

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/templates"
	"github.com/Woeter69/ini/internal/ui"
)

func init() {
	Register("cpp", &CppHandler{})
}

// CppHandler scaffolds C++ projects with g++ and a Makefile.
type CppHandler struct{}

func (c *CppHandler) Name() string {
	return "C++"
}

// SupportedTypes declares which global taxonomy categories C++ supports
func (c *CppHandler) SupportedTypes() []string {
	return []string{"basic", "app", "cli", "embedded", "os", "network", "data", "math"}
}

func (c *CppHandler) Validate() error {
	return nil
}

func (c *CppHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory + src/ and include/
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	srcDir := filepath.Join(projectDir, "src")
	includeDir := filepath.Join(projectDir, "include")
	os.MkdirAll(srcDir, 0o755)
	os.MkdirAll(includeDir, 0o755)

	// Determine type path for template
	typeDir := config.Type
	if typeDir == "" || typeDir == "basic" || typeDir == "app" {
		typeDir = "basic"
	}
	mainTmplPath := fmt.Sprintf("cpp/%s/main.cpp.tmpl", typeDir)
	makeTmplPath := fmt.Sprintf("cpp/%s/Makefile.tmpl", typeDir)

	// 2. Create src/main.cpp
	fmt.Printf("  %s Creating src/main.cpp...\n", ui.Arrow)
	if err := c.processTemplate(config, mainTmplPath, filepath.Join(srcDir, "main.cpp")); err != nil {
		return err
	}
	fmt.Printf("  %s src/main.cpp created\n", ui.CheckMark)

	// 3. Create Makefile
	fmt.Printf("  %s Creating Makefile...\n", ui.Arrow)
	if err := c.processTemplate(config, makeTmplPath, filepath.Join(projectDir, "Makefile")); err != nil {
		return err
	}
	fmt.Printf("  %s Makefile created\n", ui.CheckMark)

	// 4. Create .gitignore (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 5. Create README.md (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 6. Initialize git repo (if --git flag is set)
	if config.Git {
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	}

	// Print summary
	fmt.Println()
	relPath, _ := filepath.Rel(".", projectDir)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}

	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your C++ %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  make && make run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}

func (c *CppHandler) processTemplate(config ProjectConfig, tmplPath, destPath string) error {
	content, err := templates.FS.ReadFile(tmplPath)
	if err != nil {
		// Fallback for types that might not have a specific Makefile (use basic)
		if strings.HasSuffix(tmplPath, "Makefile.tmpl") {
			tmplPath = "cpp/basic/Makefile.tmpl"
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
