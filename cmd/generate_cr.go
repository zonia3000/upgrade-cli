package cmd

import (
	imagesettype "upgrade-cli/flag/image_set_type"
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

		imageSetType, _ := cmd.Flags().GetString(imageSetTypeFlag)

		operatorMode, _ := cmd.Flags().GetString(operatorModeFlag)
		olm, err := isOlm(operatorMode)
		if err != nil {
			return err
		}

		needsFix := service.AdaptImagesOverride(entandoApp, imagesettype.ImageSetType(imageSetType), olm)

		fileName, _ := cmd.Flags().GetString(outputFlag)
		return service.GenerateCustomResource(fileName, entandoApp, needsFix)
	},
}

func init() {
	rootCmd.AddCommand(generateCRCmd)

	generateCRCmd.Flags().StringP(outputFlag, "o", "", "path to CR file")
}
