package service

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	flag "upgrade-cli/flag"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/google/go-containerregistry/pkg/crane"
)

var craneDigest = crane.Digest

const (
	defaultRegistry              = "registry.hub.docker.com"
	defaultOrganization          = "entando"
	defaultAppBuilderImage       = "app-builder"
	defaultDeAppEapImage         = "entando-de-app-eap"
	defaultDeAppWildflyImage     = "entando-de-app-wildfly"
	defaultComponentManagerImage = "entando-component-manager"
	defaultKeycloakImage         = "entando-keycloak"
	defaultRedHatSsoImage        = "entando-redhat-sso"
)

// Convert the format of the images provided by the user to full URL format
func AdaptImagesOverride(entandoAppV2 *v1alpha1.EntandoAppV2, installationType flag.InstallationType, olm bool) error {

	defaultDeAppImage := getDefaultDeAppImage(installationType)
	defaultKeycloakImage := getDefaultKeycloakImage(installationType)

	err := adaptImageOverride(&entandoAppV2.Spec.AppBuilder.ImageOverride, defaultAppBuilderImage, olm)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.DeApp.ImageOverride, defaultDeAppImage, olm)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.ComponentManager.ImageOverride, defaultComponentManagerImage, olm)
	if err != nil {
		return err
	}
	err = adaptImageOverride(&entandoAppV2.Spec.Keycloak.ImageOverride, defaultKeycloakImage, olm)
	if err != nil {
		return err
	}

	checkInstallationTypeImagesMismatch(entandoAppV2.Spec.DeApp.ImageOverride, defaultDeAppImage, installationType)
	checkInstallationTypeImagesMismatch(entandoAppV2.Spec.Keycloak.ImageOverride, defaultKeycloakImage, installationType)

	return nil
}

func getDefaultDeAppImage(installationType flag.InstallationType) string {
	if installationType == flag.RedhatCertified {
		return defaultDeAppEapImage
	}
	return defaultDeAppWildflyImage
}

func getDefaultKeycloakImage(installationType flag.InstallationType) string {
	if installationType == flag.RedhatCertified {
		return defaultRedHatSsoImage
	}
	return defaultKeycloakImage
}

func adaptImageOverride(imageOverride *string, defaultImage string, olm bool) error {
	if *imageOverride != "" {
		if !strings.Contains(*imageOverride, ":") && !strings.Contains(*imageOverride, "/") {
			// only the tag was provided
			*imageOverride = fmt.Sprintf("%s/%s/%s:%s", defaultRegistry, defaultOrganization, defaultImage, *imageOverride)
		} else if !containsRegistry(*imageOverride) {
			*imageOverride = fmt.Sprintf("%s/%s", defaultRegistry, *imageOverride)
		}

		if olm {
			return replaceTagsWithDigests(imageOverride)
		}
	}
	return nil
}

func containsRegistry(image string) bool {
	re := regexp.MustCompile(`^[\w-\.]+\.[\w-\.]+\/[\w-\/@\.:]+$`)
	return re.MatchString(image)
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

// in case of inconsistencies between the provided images and the selected installation type the user is warned
func checkInstallationTypeImagesMismatch(image, expectedRepo string, installationType flag.InstallationType) {

	// the check is performed only when using official Entando images
	if strings.HasPrefix(image, defaultRegistry+"/"+defaultOrganization+"/") {
		re := regexp.MustCompile(`^.+/(.+):.+$`)
		matches := re.FindStringSubmatch(image)
		if len(matches) == 2 {
			providedRepo := matches[1]
			if providedRepo != expectedRepo {
				fmt.Fprintf(os.Stderr, "WARNING: installationType is set to %s but the repository %s was provided. Expected repository should be %s\n", installationType, providedRepo, expectedRepo)
			}
		}
	}
}
