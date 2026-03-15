package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Woeter69/ini/internal/handler"
	"github.com/Woeter69/ini/internal/taxonomy"
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

// makeLangCmd creates a cobra.Command for a language with --git, --type flags, interactive prompt, etc.
func makeLangCmd(cfg langCmdConfig) *cobra.Command {
	var gitFlag bool
	var typeFlag string
	var frameworkFlag string
	var variantFlag string

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

			// Determine supported types
			supported := []string{"basic"}
			if th, ok := h.(handler.TypedHandler); ok {
				supported = th.SupportedTypes()
			}

			// If typeFlag was passed manually, validate it
			if typeFlag != "" {
				typeFlag = taxonomy.Canonical(typeFlag)

				if !taxonomy.IsValid(typeFlag) {
					fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), fmt.Sprintf("invalid project type %q.", typeFlag))
					os.Exit(1)
				}

				supports := false
				for _, st := range supported {
					if taxonomy.Canonical(st) == typeFlag {
						supports = true
						break
					}
				}
				if !supports {
					fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), fmt.Sprintf("%s does not support building %q projects.", cfg.DisplayName, taxonomy.GetName(typeFlag)))
					os.Exit(1)
				}
			}

			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else {
				// Interactive Name Prompt
				nameGrp := huh.NewGroup(
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
				)
				if err := huh.NewForm(nameGrp).Run(); err != nil {
					fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
					os.Exit(1)
				}
			}

			// Interactive Type Prompt (if not passed as flag and handler supports multiple types)
			if typeFlag == "" {
				if len(supported) > 1 || (len(supported) == 1 && supported[0] != "basic") {
					var options []huh.Option[string]
					for _, st := range supported {
						options = append(options, huh.NewOption[string](taxonomy.GetName(st), st))
					}

					typeGrp := huh.NewGroup(
						huh.NewSelect[string]().
							Title(fmt.Sprintf("What kind of %s project are you building?", cfg.DisplayName)).
							Options(options...).
							Value(&typeFlag),
					)

					if err := huh.NewForm(typeGrp).Run(); err != nil {
						fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
						os.Exit(1)
					}
				} else {
					// Fallback implicitly
					typeFlag = "basic"
				}
			}

			// 3. Variant (JS vs TS) - Detect from alias first
			isJS := cfg.Lang == "bun" || cfg.Lang == "js" || cfg.Lang == "ts" || cfg.Lang == "node"
			if isJS && variantFlag == "" {
				called := cmd.CalledAs()
				if called == "js" || called == "javascript" {
					variantFlag = "js"
				} else if called == "ts" || called == "typescript" {
					variantFlag = "ts"
				}

				// If still empty and in interactive mode, ask
				if variantFlag == "" {
					langGrp := huh.NewGroup(
						huh.NewSelect[string]().
							Title("Choose Language Variant").
							Options(
								huh.NewOption("TypeScript (TSX/TS)", "ts"),
								huh.NewOption("JavaScript (JSX/JS)", "js"),
							).
							Value(&variantFlag),
					)
					if err := huh.NewForm(langGrp).Run(); err != nil {
						fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
						os.Exit(1)
					}
				}
			}

			// Special case for JS/TS Web Frameworks - Multi-stage Picker
			var category string
			if isJS && typeFlag == "web" && frameworkFlag == "" {
				// 1. Language Variant (Already handled above for JS/TS generally)

				// 2. Category
				catGrp := huh.NewGroup(
					huh.NewSelect[string]().
						Title("Frontend or Backend?").
						Options(
							huh.NewOption("Frontend (Next, React, Vue, etc.)", "frontend"),
							huh.NewOption("Backend (Express.js)", "backend"),
						).
						Value(&category),
				)
				if err := huh.NewForm(catGrp).Run(); err != nil {
					fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
					os.Exit(1)
				}

				// 3. Framework
				if category == "backend" {
					frameworkFlag = "express"
				} else {
					fwGrp := huh.NewGroup(
						huh.NewSelect[string]().
							Title("Which Framework/Library?").
							Options(
								huh.NewOption("Next.js (Recommended)", "next"),
								huh.NewOption("React", "react"),
								huh.NewOption("Vue.js", "vue"),
								huh.NewOption("Svelte", "svelte"),
								huh.NewOption("Solid", "solid"),
								huh.NewOption("Angular", "angular"),
								huh.NewOption("Vanilla (HTML/CSS)", "vanilla"),
							).
							Value(&frameworkFlag),
					)
					if err := huh.NewForm(fwGrp).Run(); err != nil {
						fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Cancelled."))
						os.Exit(1)
					}
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

			fmt.Printf("%s Initializing %s %s project %s\n\n",
				ui.TitleStyle.Render("⚡"),
				cfg.DisplayName,
				taxonomy.GetName(typeFlag),
				ui.TitleStyle.Render(fmt.Sprintf("%q", projectName)))

			config := handler.ProjectConfig{
				Name:      projectName,
				Path:      projectPath,
				Language:  cfg.Lang,
				Type:      typeFlag,
				Framework: frameworkFlag,
				Variant:   variantFlag,
				Git:       gitFlag,
			}

			if err := h.Init(config); err != nil {
				fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&gitFlag, "git", false, "Initialize a git repository with 'main' branch")
	cmd.Flags().StringVarP(&typeFlag, "type", "t", "", "Type of project to scaffold (e.g. web, game, ai)")
	cmd.Flags().StringVarP(&frameworkFlag, "framework", "f", "", "Sub-framework for JS/TS projects (e.g. react, next, express)")
	cmd.Flags().StringVarP(&variantFlag, "variant", "v", "", "Language variant for JS/TS projects (js or ts)")
	return cmd
}
