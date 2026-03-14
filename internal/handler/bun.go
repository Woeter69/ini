package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("bun", &BunHandler{}) }

type BunHandler struct{}

func (b *BunHandler) Name() string { return "JavaScript/TypeScript (Bun)" }
func (b *BunHandler) Validate() error {
	if _, err := exec.LookPath("bun"); err != nil {
		return fmt.Errorf("bun is not installed or not in PATH.\n  Install it: curl -fsSL https://bun.sh/install | bash")
	}
	return nil
}

func (b *BunHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Bun project...\n", ui.Arrow)
	bunInit := exec.Command("bun", "init", "-y") // -y bypasses prompts
	bunInit.Dir = config.Path
	bunInit.Stdout = nil
	bunInit.Stderr = os.Stderr
	if err := bunInit.Run(); err != nil {
		return fmt.Errorf("failed to init bun project: %w", err)
	}
	fmt.Printf("  %s bun project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	
	// Bun init might create git, but let's control it depending on the flag
	if config.Git {
		// we can just re-init to be safe, or leave it. Scaffold InitGit handles it gracefully
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your JavaScript/TypeScript project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  bun run index.ts\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
