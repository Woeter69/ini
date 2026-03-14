package cmd

import (
	"fmt"
	"os"

	"github.com/Woeter69/ini/internal/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ini",
	Short: "Blazing fast project initializer",
	Long:  "INI — A blazing fast CLI tool that scaffolds projects for any programming language.",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner()
		cmd.Help()
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
