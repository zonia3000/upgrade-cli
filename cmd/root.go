package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	flag "upgrade-cli/flag"
	"upgrade-cli/service"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/spf13/cobra"
)

const (
	versionFlag          = "version"
	latestFlag           = "latest"
	installationTypeFlag = "installationType"
	olmFlag              = "olm"
)

var componentFlags = map[string]string{
	"image-de-app":            "DeApp",
	"image-app-builder":       "AppBuilder",
	"image-component-manager": "ComponentManager",
	"image-keycloak":          "Keycloak",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "upgrade-cli",
	Short: "Entando Upgrade CLI",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP(versionFlag, "v", "", "Entando version")
	rootCmd.PersistentFlags().Bool(latestFlag, false, "Automatically select the latest version from entando-releases repository")
	rootCmd.MarkFlagsMutuallyExclusive(versionFlag, latestFlag)

	installationTypeFlagValue := flag.GetInstallationTypeFlag()
	installationTypeFlagUsage := "Set specific images for DeApp or Keycloak. Possible values: " + strings.Join(flag.GetInstallationTypeValues(), ", ")
	rootCmd.PersistentFlags().VarP(installationTypeFlagValue, installationTypeFlag, "t", installationTypeFlagUsage)

	olmFlagValue := flag.GetBoolOrAutoFlag()
	olmFlagUsage := "Generate CR for an OLM installation. Possible values: " + strings.Join(flag.GetBoolOrAutoValues(), ", ")
	rootCmd.PersistentFlags().Var(olmFlagValue, olmFlag, olmFlagUsage)

	// Global component flags
	for componentFlag, componentName := range componentFlags {
		rootCmd.PersistentFlags().String(componentFlag, "", "Image override for "+componentName)
	}
}

func ParseEntandoAppFromCmd(cmd *cobra.Command) (*v1alpha1.EntandoAppV2, error) {

	var version string
	if latest, _ := cmd.Flags().GetBool(latestFlag); latest {
		version = service.GetLatestVersion()
	} else {
		version, _ = cmd.Flags().GetString(versionFlag)
	}

	entandoApp := v1alpha1.EntandoAppV2{}
	entandoApp.Spec.Version = version

	err := parseComponentFlag(cmd, "image-de-app", &entandoApp.Spec.DeApp.ImageOverride)
	if err != nil {
		return nil, err
	}
	err = parseComponentFlag(cmd, "image-app-builder", &entandoApp.Spec.AppBuilder.ImageOverride)
	if err != nil {
		return nil, err
	}
	err = parseComponentFlag(cmd, "image-component-manager", &entandoApp.Spec.ComponentManager.ImageOverride)
	if err != nil {
		return nil, err
	}
	err = parseComponentFlag(cmd, "image-keycloak", &entandoApp.Spec.Keycloak.ImageOverride)
	if err != nil {
		return nil, err
	}

	return &entandoApp, nil
}

func isOlm(olmFlagValue string) (bool, error) {
	switch flag.BoolOrAutoType(olmFlagValue) {
	case flag.True:
		return true, nil
	case flag.False:
		return false, nil
	}
	return false, fmt.Errorf("automatic detection of OLM not implemented yet")
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
