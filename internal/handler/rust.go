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
	Register("rust", &RustHandler{})
}

// RustHandler scaffolds Rust projects using cargo.
type RustHandler struct{}

func (r *RustHandler) Name() string {
	return "Rust"
}

// SupportedTypes declares which global taxonomy categories Rust supports
func (r *RustHandler) SupportedTypes() []string {
	return []string{
		"basic", "web", "script", "game", "network", "os", "db",
		"security", "graphics", "web3", "lang",
	}
}

func (r *RustHandler) Validate() error {
	_, err := exec.LookPath("cargo")
	if err != nil {
		return fmt.Errorf("cargo is not installed or not in PATH.\n  Install it: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh")
	}
	return nil
}

func (r *RustHandler) Init(config ProjectConfig) error {
	projectDir := config.Path

	// 1. Create project directory
	if err := scaffold.CreateDir(projectDir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 2. Run cargo init inside the directory
	fmt.Printf("  %s Initializing Rust project with cargo...\n", ui.Arrow)
	cargoInit := exec.Command("cargo", "init", "--name", config.Name)
	cargoInit.Dir = projectDir
	cargoInit.Stdout = nil
	cargoInit.Stderr = os.Stderr
	if err := cargoInit.Run(); err != nil {
		return fmt.Errorf("failed to init cargo project: %w", err)
	}
	fmt.Printf("  %s Cargo project initialized\n", ui.CheckMark)

	// Determine type path for template
	templatePath := "rust/basic/main.rs.tmpl" // fallback
	deps := []string{}

	switch config.Type {
	case "web":
		templatePath = "rust/web/main.rs.tmpl"
		deps = append(deps, "axum", "tokio@1", "tokio-features=macros,rt-multi-thread")
	case "script":
		templatePath = "rust/script/main.rs.tmpl"
		deps = append(deps, "clap@4", "clap-features=derive")
	case "game":
		templatePath = "rust/game/main.rs.tmpl"
		deps = append(deps, "bevy")
	case "network":
		templatePath = "rust/network/main.rs.tmpl"
		deps = append(deps, "tokio@1", "tokio-features=full")
	case "os":
		templatePath = "rust/os/main.rs.tmpl"
		deps = append(deps, "sysinfo")
	case "db":
		templatePath = "rust/db/main.rs.tmpl"
		deps = append(deps, "rusqlite@0.31", "rusqlite-features=bundled")
	case "security":
		templatePath = "rust/security/main.rs.tmpl"
		deps = append(deps, "sha2", "rand")
	case "graphics":
		templatePath = "rust/graphics/main.rs.tmpl"
		deps = append(deps, "winit")
	case "web3":
		templatePath = "rust/web3/main.rs.tmpl"
		deps = append(deps, "ethers", "tokio@1", "tokio-features=macros,rt-multi-thread")
	case "lang":
		templatePath = "rust/lang/main.rs.tmpl"
		deps = append(deps, "syn", "quote")
	case "basic":
		templatePath = "rust/basic/main.rs.tmpl"
	}

	// 3. Create src/main.rs from embedded template
	fmt.Printf("  %s Creating main.rs...\n", ui.Arrow)
	tmplContent, err := templates.FS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	t, err := template.New("main").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(filepath.Join(projectDir, "src", "main.rs"), buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write main.rs: %w", err)
	}
	fmt.Printf("  %s main.rs created\n", ui.CheckMark)

	// 4. Resolve dependencies if any
	if len(deps) > 0 {
		fmt.Printf("  %s Fetching dependencies...\n", ui.Arrow)
		for _, dep := range deps {
			var addArgs []string
			addArgs = append(addArgs, "add")
			if strings.Contains(dep, "-features=") {
				parts := strings.Split(dep, "-features=")
				addArgs = append(addArgs, parts[0], "--features", parts[1])
			} else {
				addArgs = append(addArgs, dep)
			}

			cargoAdd := exec.Command("cargo", addArgs...)
			cargoAdd.Dir = projectDir
			cargoAdd.Stdout = nil
			cargoAdd.Stderr = os.Stderr
			if err := cargoAdd.Run(); err != nil {
				return fmt.Errorf("failed to add dependency '%s': %w", dep, err)
			}
		}
		fmt.Printf("  %s Dependencies resolved\n", ui.CheckMark)
	}

	// 5. Overwrite .gitignore with our comprehensive one (shared)
	if err := scaffold.WriteGitignore(projectDir, config.Language); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// 6. Overwrite README.md with our template (shared)
	if err := scaffold.WriteReadme(projectDir, config.Name, config.Language); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// 7. Initialize git repo (if --git flag is set)
	// Note: cargo init creates a .git by default, so we remove it first
	// and let our scaffold handle it with -b main
	if config.Git {
		// cargo init may have already created a git repo, remove it so we get -b main
		os.RemoveAll(filepath.Join(projectDir, ".git"))
		if err := scaffold.InitGit(projectDir); err != nil {
			return err
		}
	} else {
		// Remove the git repo cargo created since user didn't ask for it
		os.RemoveAll(filepath.Join(projectDir, ".git"))
	}

	// Print summary
	fmt.Println()
	relPath, _ := filepath.Rel(".", projectDir)
	if relPath == "" || relPath == "." {
		relPath = config.Name
	}
	summary := strings.Builder{}
	summary.WriteString(ui.SuccessStyle.Render(fmt.Sprintf("🚀 Your Rust %s project is ready!", config.Type)))
	summary.WriteString("\n\n")
	summary.WriteString(fmt.Sprintf("  cd %s\n", relPath))
	summary.WriteString("  cargo run\n")

	fmt.Println(ui.SummaryBox.Render(summary.String()))
	fmt.Println()

	return nil
}
