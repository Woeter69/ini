package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
		var langOptions []huh.Option[string]

		// Dynamically collect all language subcommands
		for _, c := range cmd.Commands() {
			if c.Name() == "help" || c.Name() == "completion" || c.Hidden {
				continue
			}

			// We use the first word of Short as a proxy for the formal name, 
			// or just the capitalize command name.
			displayName := strings.Title(c.Name())
			if strings.Contains(c.Short, "Initialize a new") {
				// Try to extract the display name from the Short description if possible
				parts := strings.Split(c.Short, " ")
				if len(parts) >= 4 {
					displayName = parts[3]
				}
			}

			langOptions = append(langOptions, huh.NewOption(displayName, c.Name()))
		}

		// Sort options alphabetically by name
		sort.Slice(langOptions, func(i, j int) bool {
			return langOptions[i].Key < langOptions[j].Key
		})

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose your project's primary language").
					Description("ini will help you scaffold the perfect boilerplate.").
					Height(15). // Add height for longer list
					Options(langOptions...).
					Value(&selectedLang),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("Selection cancelled.")
			return
		}

		// Find and execute the selected language subcommand
		subCmd, _, err := cmd.Find([]string{selectedLang})
		if err != nil {
			fmt.Printf("Error finding subcommand %q: %v\n", selectedLang, err)
			return
		}

		if subCmd != cmd && subCmd.Run != nil {
			// Trigger the subcommand's Run function without any args,
			// forcing it to enter its own interactive mode (Name/Type/etc.)
			subCmd.Run(subCmd, []string{})
		}
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
