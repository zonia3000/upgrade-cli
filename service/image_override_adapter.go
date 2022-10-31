package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/google/go-containerregistry/pkg/crane"
)

const defaultRegistry = "registry.hub.docker.com"

var craneDigest = crane.Digest

// Convert the format of the images provided by the user to full URL format
func AdaptImagesOverride(entandoAppV2 *v1alpha1.EntandoAppV2) error {

	err := adaptImageOverride(&entandoAppV2.Spec.AppBuilder.ImageOverride)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.DeApp.ImageOverride)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.ComponentManager.ImageOverride)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.Keycloak.ImageOverride)
	if err != nil {
		return err
	}

	return nil
}

func adaptImageOverride(imageOverride *string) error {
	if *imageOverride != "" {
		addMissingRegistry(imageOverride)
		return replaceTagsWithDigests(imageOverride)
	}
	return nil
}

func addMissingRegistry(imageOverride *string) {
	re := regexp.MustCompile(`^[\w-\.]+\.[\w-\.]+\/[\w-\/@\.:]+$`)
	if !re.MatchString(*imageOverride) {
		*imageOverride = defaultRegistry + "/" + *imageOverride
	}
}

func replaceTagsWithDigests(imageOverride *string) error {
	if !strings.Contains(*imageOverride, "@sha256:") {
		digest, err := craneDigest(*imageOverride)

		if err != nil {
			return fmt.Errorf("unable to retrieve digest for image %s.\n%s", *imageOverride, err.Error())
		}

		*imageOverride = strings.Split(*imageOverride, ":")[0] + "@" + digest
	}
	return nil
}
