package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Woeter69/ini/internal/handler"
	_ "github.com/Woeter69/ini/internal/handler" // Ensure all handlers are registered
)

func TestAllHandlers(t *testing.T) {
	// Create a temporary directory for test projects
	tmpDir, err := os.MkdirTemp("", "ini-tests-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Get all registered languages from the registry
	// Since registry is unexported, we might need to export it or use a helper.
	// However, we can use the manual list or auditlangs.go logic if we had one.
	// For now, let's list them based on our known 39 languages.
	languages := []string{
		"go", "python", "rust", "bun", "java", "kotlin", "swift", "c", "cpp",
		"csharp", "zig", "nim", "fsharp", "scala", "ruby", "haskell", "lua",
		"perl", "php", "fortran", "cobol", "d", "v", "crystal", "ada",
		"julia", "r", "pascal", "clojure", "groovy", "ocaml", "erlang", "elixir",
		"objc", "assembly", "shell",
	}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			h, err := handler.Get(lang)
			if err != nil {
				t.Fatalf("Failed to get handler for %s: %v", lang, err)
			}

			// Skip if toolchain is missing
			if err := h.Validate(); err != nil {
				t.Skipf("Skipping %s: %v", lang, err)
			}

			// Test Basic Scaffolding
			testScaffold(t, h, lang, "basic", tmpDir)

			// Test CLI Scaffolding if supported
			if th, ok := h.(handler.TypedHandler); ok {
				for _, pType := range th.SupportedTypes() {
					if pType == "basic" { continue }
					t.Run(pType, func(t *testing.T) {
						testScaffold(t, h, lang, pType, tmpDir)
					})
				}
			}
		})
	}
}

func testScaffold(t *testing.T, h handler.Handler, lang, pType, tmpParent string) {
	projectName := fmt.Sprintf("%s-%s-project", lang, pType)
	projectPath := filepath.Join(tmpParent, projectName)

	config := handler.ProjectConfig{
		Name:     projectName,
		Path:     projectPath,
		Language: lang,
		Type:     pType,
		Git:      false,
	}

	if err := h.Init(config); err != nil {
		t.Fatalf("Scaffolding failed for %s (%s): %v", lang, pType, err)
	}

	// Verify existence of the project directory
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Errorf("Project directory not created for %s (%s)", lang, pType)
	}

	// Verify some common files
	filesToVerify := []string{"README.md", ".gitignore"}
	for _, f := range filesToVerify {
		if _, err := os.Stat(filepath.Join(projectPath, f)); os.IsNotExist(err) {
			t.Errorf("File %s missing for %s (%s)", f, lang, pType)
		}
	}
}
