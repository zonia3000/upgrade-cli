package cmd

import (
	"upgrade-cli/service"

	"github.com/spf13/cobra"
)

const (
	outputFlag = "output"
)

var generateCRCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate EntandoAppV2 CR file",
	PreRun: func(cmd *cobra.Command, args []string) {
		latest, _ := cmd.Flags().GetBool(latestFlag)
		if !latest {
			cmd.MarkFlagRequired(versionFlag)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prevent showing usage message when error happens in RunE func
		cmd.SilenceUsage = true

		entandoApp, err := ParseEntandoAppFromCmd(cmd)
		if err != nil {
			return err
		}

		fileName, _ := cmd.Flags().GetString(outputFlag)
		service.GenerateCustomResource(fileName, entandoApp)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCRCmd)

	generateCRCmd.Flags().StringP(outputFlag, "o", "", "path to CR file")
}
