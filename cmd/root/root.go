package root

import (
	"os"

	"upgrade-cli/cmd/generate"
	"upgrade-cli/cmd/upgrade"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "upgrade-cli",
	Short: "Entando Upgrade CLI",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.AddCommand(generate.GenerateCRCmd)
	RootCmd.AddCommand(upgrade.UpgradeCmd)
}
