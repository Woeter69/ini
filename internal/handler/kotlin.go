package handler

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Woeter69/ini/internal/scaffold"
	"github.com/Woeter69/ini/internal/ui"
)

func init() { Register("kotlin", &KotlinHandler{}) }

type KotlinHandler struct{}

func (k *KotlinHandler) Name() string    { return "Kotlin" }
func (k *KotlinHandler) Validate() error {
	if _, err := exec.LookPath("gradle"); err != nil {
		return fmt.Errorf("gradle is not installed or not in PATH.\n  Install it: https://gradle.org/install/")
	}
	return nil
}

func (k *KotlinHandler) Init(config ProjectConfig) error {
	if err := scaffold.CreateDir(config.Path); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	fmt.Printf("  %s Initializing Kotlin project with Gradle...\n", ui.Arrow)
	gradleInit := exec.Command("gradle", "init", "--type", "kotlin-application", "--dsl", "kotlin", "--test-framework", "kotlintest", "--project-name", config.Name, "--package", config.Name, "--no-split-project", "--java-version", "21")
	gradleInit.Dir = config.Path
	gradleInit.Stdout = nil
	gradleInit.Stderr = os.Stderr
	gradleInit.Stdin = nil
	if err := gradleInit.Run(); err != nil {
		return fmt.Errorf("failed to init gradle project: %w", err)
	}
	fmt.Printf("  %s Gradle Kotlin project initialized\n", ui.CheckMark)

	if err := scaffold.WriteGitignore(config.Path, config.Language); err != nil {
		return err
	}
	if err := scaffold.WriteReadme(config.Path, config.Name, config.Language); err != nil {
		return err
	}
	if config.Git {
		os.RemoveAll(config.Path + "/.git")
		if err := scaffold.InitGit(config.Path); err != nil { return err }
	} else {
		os.RemoveAll(config.Path + "/.git")
	}

	fmt.Println()
	s := strings.Builder{}
	s.WriteString(ui.SuccessStyle.Render("🚀 Your Kotlin project is ready!"))
	s.WriteString("\n\n")
	s.WriteString(fmt.Sprintf("  cd %s\n", config.Name))
	s.WriteString("  gradle run\n")
	fmt.Println(ui.SummaryBox.Render(s.String()))
	fmt.Println()
	return nil
}
