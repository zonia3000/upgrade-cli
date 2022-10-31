package cmd

import (
	"fmt"
	"regexp"
	"upgrade-cli/service"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/spf13/cobra"
)

const (
	fileFlag    = "file"
	versionFlag = "version"
	latestFlag  = "latest"
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

		entandoAppV2 := v1alpha1.EntandoAppV2{}
		entandoAppV2.Spec.Version = version

		err := parseComponentFlag(cmd, "de-app", &entandoAppV2.Spec.DeApp.ImageOverride)
		if err != nil {
			return err
		}
		err = parseComponentFlag(cmd, "app-builder", &entandoAppV2.Spec.AppBuilder.ImageOverride)
		if err != nil {
			return err
		}
		err = parseComponentFlag(cmd, "component-manager", &entandoAppV2.Spec.ComponentManager.ImageOverride)
		if err != nil {
			return err
		}
		err = parseComponentFlag(cmd, "keycloak", &entandoAppV2.Spec.Keycloak.ImageOverride)
		if err != nil {
			return err
		}

		err = service.AdaptImagesOverride(&entandoAppV2)
		if err != nil {
			return err
		}

		service.GenerateCustomResource(file, &entandoAppV2)

		return nil
	},
}

func parseComponentFlag(cmd *cobra.Command, componentFlag string, imageOverride *string) error {
	componentImage, _ := cmd.Flags().GetString(componentFlag)

	if componentImage != "" {
		re := regexp.MustCompile(`^[\w-\/\.@]*:?[\w-\.]+$`)

		if !re.MatchString(componentImage) {
			return fmt.Errorf("invalid format for image override flag '%s'. It should be <image>:<tag> or <tag>", componentImage)
		}

		*imageOverride = componentImage
	}

	return nil
}

func init() {
	rootCmd.AddCommand(generateCRCmd)

	generateCRCmd.Flags().StringP(fileFlag, "f", "", "path to CR file")
	generateCRCmd.MarkFlagRequired(fileFlag)

	generateCRCmd.Flags().StringP(versionFlag, "v", "", "Entando version")
	generateCRCmd.Flags().Bool(latestFlag, false, "Automatically select the latest version from entando-releases repository")
	generateCRCmd.MarkFlagsMutuallyExclusive(versionFlag, latestFlag)

	addComponentFlag("de-app", "DeApp")
	addComponentFlag("app-builder", "AppBuilder")
	addComponentFlag("component-manager", "ComponentManager")
	addComponentFlag("keycloak", "Keycloak")
}

func addComponentFlag(componentFlag, componentName string) {
	generateCRCmd.Flags().String(componentFlag, "", "Image override for "+componentName)
}
