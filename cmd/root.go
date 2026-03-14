package cmd

import (
	"fmt"
	"os"

	"github.com/Woeter69/ini/internal/ui"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ini",
	Short: "Blazing fast project initializer",
	Long:  "INI — A blazing fast CLI tool that scaffolds projects for any programming language.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		ui.PrintBanner()
		
		var selectedLang string
		langOptions := []huh.Option[string]{
			huh.NewOption("Go", "go"),
			huh.NewOption("Python", "python"),
			huh.NewOption("Rust", "rust"),
			huh.NewOption("JavaScript/TypeScript (Bun)", "bun"),
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose your project's primary language").
					Description("ini will help you scaffold the perfect boilerplate.").
					Options(langOptions...).
					Value(&selectedLang),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("Selection cancelled.")
			return
		}

		// Execute the selected language subcommand
		cmd.SetArgs([]string{selectedLang})
		cmd.Parent().Execute() // This avoids direct self-reference issues if possible
	},
	Version: ui.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, ui.ErrorStyle.Render("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("ini v%s\n", ui.Version))
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
