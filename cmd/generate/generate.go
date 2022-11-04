package generate

import (
	"fmt"
	"regexp"
	"strings"
	"upgrade-cli/flag/component"
	imagesettype "upgrade-cli/flag/image_set_type"
	operatormode "upgrade-cli/flag/operator_mode"
	"upgrade-cli/service"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/spf13/cobra"
)

const (
	// Flags shared with the upgrade command
	VersionFlag       = "version"
	LatestVersionFlag = "latest-version"
	ImageSetTypeFlag  = "image-set-type"
	OperatorModeFlag  = "operator-mode"

	// Flag specific of the generate command
	outputFlag = "output"
)

var GenerateCRCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate EntandoAppV2 CR file",
	PreRun: func(cmd *cobra.Command, args []string) {
		latest, _ := cmd.Flags().GetBool(LatestVersionFlag)
		if !latest {
			cmd.MarkFlagRequired(VersionFlag)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prevent showing usage message when error happens in RunE func
		cmd.SilenceUsage = true

		entandoApp, err := ParseEntandoAppFromCmd(cmd)
		if err != nil {
			return err
		}

		imageSetType, _ := cmd.Flags().GetString(ImageSetTypeFlag)

		operatorMode, _ := cmd.Flags().GetString(OperatorModeFlag)
		olm, err := IsOlm(operatorMode)
		if err != nil {
			return err
		}

		needsFix := service.AdaptImagesOverride(entandoApp, imagesettype.ImageSetType(imageSetType), olm)

		fileName, _ := cmd.Flags().GetString(outputFlag)
		return service.GenerateCustomResource(fileName, entandoApp, needsFix)
	},
}

func init() {
	AddCRFlags(GenerateCRCmd)

	GenerateCRCmd.Flags().StringP(outputFlag, "o", "", "path to CR file")
}

func ParseEntandoAppFromCmd(cmd *cobra.Command) (*v1alpha1.EntandoAppV2, error) {

	var version string
	if latest, _ := cmd.Flags().GetBool(LatestVersionFlag); latest {
		version = service.GetLatestVersion()
	} else {
		version, _ = cmd.Flags().GetString(VersionFlag)
	}

	entandoApp := v1alpha1.EntandoAppV2{}
	entandoApp.Spec.Version = version

	for _, componentFlag := range component.ComponentFlags {
		err := parseComponentFlag(cmd, componentFlag, &entandoApp)
		if err != nil {
			return nil, err
		}
	}

	return &entandoApp, nil
}

func IsOlm(operatorModeFlagValue string) (bool, error) {
	switch operatormode.OperatorMode(operatorModeFlagValue) {
	case operatormode.OLM:
		return true, nil
	case operatormode.Plain:
		return false, nil
	}
	return false, fmt.Errorf("automatic detection of OLM not implemented yet")
}

func parseComponentFlag(cmd *cobra.Command, componentFlag component.ComponentFlag, entandoApp *v1alpha1.EntandoAppV2) error {
	componentImage, _ := cmd.Flags().GetString(componentFlag.Flag)

	if componentImage != "" {
		re := regexp.MustCompile(`^[\w-\/\.@]*:?[\w-\.]+$`)

		if !re.MatchString(componentImage) {
			return fmt.Errorf("invalid format for image override flag '%s'. It should be <image>:<tag> or <tag>", componentImage)
		}

		imageOverride := componentFlag.ImageOverrideGetter(entandoApp)
		*imageOverride = componentImage
	}

	return nil
}

func AddCRFlags(cmd *cobra.Command) {

	cmd.PersistentFlags().StringP(VersionFlag, "v", "", "Entando version")
	cmd.PersistentFlags().Bool(LatestVersionFlag, false, "Automatically select the latest version from entando-releases repository")
	cmd.MarkFlagsMutuallyExclusive(VersionFlag, LatestVersionFlag)

	imageSetTypeFlagValue := imagesettype.GetImageSetTypeFlag()
	imageSetTypeFlagUsage := "Set specific images for DeApp or Keycloak. Possible values: " + strings.Join(imagesettype.GetImageSetTypeValues(), ", ")
	cmd.PersistentFlags().VarP(imageSetTypeFlagValue, ImageSetTypeFlag, "t", imageSetTypeFlagUsage)

	operatorModeFlagValue := operatormode.GetOperatorModeFlag()
	operatorModeFlagUsage := "Generate CR for an OLM or plain installation. Possible values: " + strings.Join(operatormode.GetOperatorModeValues(), ", ")
	cmd.PersistentFlags().VarP(operatorModeFlagValue, OperatorModeFlag, "m", operatorModeFlagUsage)

	for _, componentFlag := range component.ComponentFlags {
		cmd.PersistentFlags().String(componentFlag.Flag, "", "Image override for "+componentFlag.ComponentName)
	}
}