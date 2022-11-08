package generate

import (
	"fmt"
	"regexp"
	"strings"
	imagesettype "upgrade-cli/flag/image_set_type"
	operatormode "upgrade-cli/flag/operator_mode"
	"upgrade-cli/service"
	"upgrade-cli/util/images"

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

		entandoApp, olm, err := ParseEntandoAppFromCmd(cmd)
		if err != nil {
			return err
		}

		needsFix := service.AdaptImagesOverride(entandoApp, olm)

		fileName, _ := cmd.Flags().GetString(outputFlag)
		return service.GenerateCustomResource(fileName, entandoApp, needsFix)
	},
}

func init() {
	AddCRFlags(GenerateCRCmd)

	GenerateCRCmd.Flags().StringP(outputFlag, "o", "", "path to CR file")
}

func ParseEntandoAppFromCmd(cmd *cobra.Command) (*v1alpha1.EntandoAppV2, bool, error) {

	var version string
	if latest, _ := cmd.Flags().GetBool(LatestVersionFlag); latest {
		version = service.GetLatestVersion()
	} else {
		version, _ = cmd.Flags().GetString(VersionFlag)
	}

	olm, err := isOlm(cmd)
	if err != nil {
		return nil, false, err
	}

	imageSetType := getImageSetType(cmd, olm)

	entandoApp := v1alpha1.EntandoAppV2{}
	entandoApp.Spec.Version = version
	entandoApp.Spec.ImageSetType = string(imageSetType)

	for _, imageInfo := range images.EntandoImages {
		err := parseComponentFlag(cmd, imageInfo, &entandoApp)
		if err != nil {
			return nil, false, err
		}
	}

	return &entandoApp, olm, nil
}

func isOlm(cmd *cobra.Command) (bool, error) {
	flagValue, _ := cmd.Flags().GetString(OperatorModeFlag)
	if flagValue == string(operatormode.Auto) {
		mode, err := service.GetOperatorMode()
		if err != nil {
			return false, err
		}
		return mode == operatormode.OLM, nil
	}
	return flagValue == string(operatormode.OLM), nil
}

func getImageSetType(cmd *cobra.Command, olm bool) imagesettype.ImageSetType {
	flagValue, _ := cmd.Flags().GetString(ImageSetTypeFlag)
	if flagValue == string(imagesettype.Auto) {
		if olm {
			return imagesettype.RedhatCertified
		} else {
			return imagesettype.Community
		}
	}
	return imagesettype.ImageSetType(flagValue)
}

func parseComponentFlag(cmd *cobra.Command, imageInfo images.EntandoImageInfo, entandoApp *v1alpha1.EntandoAppV2) error {
	componentImage, _ := cmd.Flags().GetString(imageInfo.ImageOverrideFlag)

	if componentImage != "" {
		re := regexp.MustCompile(`^[\w-\/\.@]*:?[\w-\.]+$`)

		if !re.MatchString(componentImage) {
			return fmt.Errorf("invalid format for image override flag '%s'. It should be <image>:<tag> or <tag>", componentImage)
		}

		imageOverride := imageInfo.GetImageOverride(entandoApp)
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

	for _, imageInfo := range images.EntandoImages {
		cmd.PersistentFlags().String(imageInfo.ImageOverrideFlag, "", "Image override for "+imageInfo.ComponentName)
	}
}
