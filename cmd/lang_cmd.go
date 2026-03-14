package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/handler"
	"github.com/Woeter69/ini/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// langCmdConfig holds config for generating a language subcommand.
type langCmdConfig struct {
	Use         string
	Aliases     []string
	Lang        string   // handler registry key
	DisplayName string
	Short       string
	Long        string
	Placeholder string
}

// makeLangCmd creates a cobra.Command for a language with --git flag, interactive prompt, etc.
func makeLangCmd(cfg langCmdConfig) *cobra.Command {
	var gitFlag bool

	cmd := &cobra.Command{
		Use:     cfg.Use + " [project-name]",
		Aliases: cfg.Aliases,
		Short:   cfg.Short,
		Long:    cfg.Long,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ui.PrintBanner()

			h, err := handler.Get(cfg.Lang)
			if err != nil {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
				os.Exit(1)
			}

			if err := h.Validate(); err != nil {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
				os.Exit(1)
			}

			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else {
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("Project name").
							Description(fmt.Sprintf("Name for your new %s project", cfg.DisplayName)).
							Placeholder(cfg.Placeholder).
							Validate(func(s string) error {
								s = strings.TrimSpace(s)
								if s == "" {
									return fmt.Errorf("project name cannot be empty")
								}
								if strings.ContainsAny(s, " /\\:*?\"<>|") {
									return fmt.Errorf("project name contains invalid characters")
								}
								return nil
							}).
							Value(&projectName),
					),
				)

				if err := form.Run(); err != nil {
					fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
					os.Exit(1)
				}
			}

			projectName = strings.TrimSpace(projectName)

			projectPath, err := filepath.Abs(projectName)
			if err != nil {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
				os.Exit(1)
			}

			if entries, err := os.ReadDir(projectPath); err == nil && len(entries) > 0 {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"),
					fmt.Sprintf("directory %q already exists and is not empty", projectName))
				os.Exit(1)
			}

			fmt.Printf("%s Initializing %s project %s\n\n",
				ui.TitleStyle.Render("⚡"),
				cfg.DisplayName,
				ui.TitleStyle.Render(fmt.Sprintf("%q", projectName)))

			config := handler.ProjectConfig{
				Name:     projectName,
				Path:     projectPath,
				Language: cfg.Lang,
				Git:      gitFlag,
			}

			if err := h.Init(config); err != nil {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&gitFlag, "git", false, "Initialize a git repository with 'main' branch")
	return cmd
}
