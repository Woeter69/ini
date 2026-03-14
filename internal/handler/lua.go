package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("lua", &LuaHandler{}) }

type LuaHandler struct{}

func (l *LuaHandler) Name() string    { return "Lua" }
func (l *LuaHandler) Validate() error { return nil }

func (l *LuaHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Creating main.lua...\n", ui.Arrow)
	mainLua := fmt.Sprintf("-- %s\n\nprint(\"Hello from %s!\")\n", config.Name, config.Name)
	if err := os.WriteFile(filepath.Join(config.Path, "main.lua"), []byte(mainLua), 0o644); err != nil {
		return err
	}
	fmt.Printf("  %s main.lua created\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil { return err }
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil { return err }
	if config.Git { if err := scaffold.InitGit(config.Path); err != nil { return err } }

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Lua project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  lua main.lua\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
