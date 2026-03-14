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

func init() { Register("clojure", &ClojureHandler{}) }

type ClojureHandler struct{}

func (c *ClojureHandler) Name() string { return "Clojure" }
func (c *ClojureHandler) Validate() error {
	if _, err := exec.LookPath("clj"); err != nil {
		return fmt.Errorf("clj (Clojure CLI) is not installed.\n  Install it: https://clojure.org/guides/install_clojure")
	}
	return nil
}

func (c *ClojureHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Clojure CLI doesn't natively scaffold an app structure like Leiningen's `lein new` does,
	// but Clojure CLI uses deps.edn. We'll set up a standard deps.edn structure.
	fmt.Printf("  %s Creating deps.edn suite...\n", ui.Arrow)

	os.MkdirAll(filepath.Join(config.Path, "src"), 0o755)

	depsEdn := `{
  :paths ["src"]
  :aliases {
    :run {:main-opts ["-m" "main"]}
  }
}`
	if err := os.WriteFile(filepath.Join(config.Path, "deps.edn"), []byte(depsEdn), 0o644); err != nil {
		return err
	}

	mainClj := fmt.Sprintf(`(ns main)

(defn -main [& args]
  (println "Hello from %s!"))
`, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "src", "main.clj"), []byte(mainClj), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s deps.edn project created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Clojure project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  clj -M:run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
