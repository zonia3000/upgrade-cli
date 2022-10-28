package cmd

import (
	"fmt"
	"regexp"
	"upgrade-cli/service"

	"github.com/spf13/cobra"
)

const (
	fileFlag    = "file"
	versionFlag = "version"
	latestFlag  = "latest"
	imageFlag   = "image"
)

var generateCRCmd = &cobra.Command{
	Use:   "generate-cr",
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

		var version string
		if latest, _ := cmd.Flags().GetBool(latestFlag); latest {
			version = service.GetLatestVersion()
		} else {
			version, _ = cmd.Flags().GetString(versionFlag)
		}

		file, _ := cmd.Flags().GetString(fileFlag)

		imagesOverrideFlags, _ := cmd.Flags().GetStringArray(imageFlag)

		imagesOverride, err := parseImagesOverride(imagesOverrideFlags)
		if err != nil {
			return err
		}

		err = service.AdaptImagesOverride(imagesOverride)
		if err != nil {
			return err
		}

		service.GenerateCustomResource(file, version, imagesOverride)

		return nil
	},
}

func parseImagesOverride(imagesOverrideFlags []string) (map[string]string, error) {
	imagesOverride := make(map[string]string)
	if len(imagesOverrideFlags) > 0 {
		for _, imageOverrideFlag := range imagesOverrideFlags {
			re := regexp.MustCompile(`^([\w-]+)=([\w-\/\.@]+:[\w-\.]+)$`)
			match := re.FindStringSubmatch(imageOverrideFlag)

			if len(match) != 3 {
				return nil, fmt.Errorf("invalid format for image override flag '%s'. It should be <component-name>=<image>:<tag>", imageOverrideFlag)
			}

			imagesOverride[match[1]] = match[2]
		}
	}
	return imagesOverride, nil
}

func init() {
	rootCmd.AddCommand(generateCRCmd)

	generateCRCmd.Flags().StringP(fileFlag, "f", "", "path to CR file")
	generateCRCmd.MarkFlagRequired(fileFlag)

	generateCRCmd.Flags().StringP(versionFlag, "v", "", "Entando version")
	generateCRCmd.Flags().Bool(latestFlag, false, "Automatically select the latest version from entando-releases repository")
	generateCRCmd.MarkFlagsMutuallyExclusive(versionFlag, latestFlag)

	generateCRCmd.Flags().StringArrayP(imageFlag, "i", []string{}, "Image override for a specific component using format <component-name>=<image>:<tag>")
}
