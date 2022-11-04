package service

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	imagesettype "upgrade-cli/flag/image_set_type"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/google/go-containerregistry/pkg/crane"
)

var CraneDigest = crane.Digest

const (
	defaultRegistry              = "registry.hub.docker.com"
	defaultOrganization          = "entando"
	defaultAppBuilderImage       = "app-builder"
	defaultDeAppEapImage         = "entando-de-app-eap"
	defaultDeAppWildflyImage     = "entando-de-app-wildfly"
	defaultComponentManagerImage = "entando-component-manager"
	defaultKeycloakImage         = "entando-keycloak"
	defaultRedHatSsoImage        = "entando-redhat-sso"

	missingDigestPlaceholder = "###-FIXME-INSERT-SHA256-###"
)

// AdaptImagesOverride converts the format of the images provided by the user to full URL format
// Returns a bool that is true in case of errors in digests retrieval. This will be used to add a comment on the YAML file
func AdaptImagesOverride(entandoAppV2 *v1alpha1.EntandoAppV2, imageSetType imagesettype.ImageSetType, olm bool) bool {

	defaultDeAppImage := getDefaultDeAppImage(imageSetType)
	defaultKeycloakImage := getDefaultKeycloakImage(imageSetType)

	digestErrors := make(map[string]error)

	adaptImageOverride(&entandoAppV2.Spec.AppBuilder.ImageOverride, defaultAppBuilderImage, olm, digestErrors)
	adaptImageOverride(&entandoAppV2.Spec.DeApp.ImageOverride, defaultDeAppImage, olm, digestErrors)
	adaptImageOverride(&entandoAppV2.Spec.ComponentManager.ImageOverride, defaultComponentManagerImage, olm, digestErrors)
	adaptImageOverride(&entandoAppV2.Spec.Keycloak.ImageOverride, defaultKeycloakImage, olm, digestErrors)

	checkInstallationTypeImagesMismatch(entandoAppV2.Spec.DeApp.ImageOverride, defaultDeAppImage, imageSetType)
	checkInstallationTypeImagesMismatch(entandoAppV2.Spec.Keycloak.ImageOverride, defaultKeycloakImage, imageSetType)

	return checkDigestErrors(digestErrors)
}

func getDefaultDeAppImage(imageSetType imagesettype.ImageSetType) string {
	if imageSetType == imagesettype.RedhatCertified {
		return defaultDeAppEapImage
	}
	return defaultDeAppWildflyImage
}

func getDefaultKeycloakImage(imageSetType imagesettype.ImageSetType) string {
	if imageSetType == imagesettype.RedhatCertified {
		return defaultRedHatSsoImage
	}
	return defaultKeycloakImage
}

func adaptImageOverride(imageOverride *string, defaultImage string, olm bool, digestErrors map[string]error) {
	if *imageOverride != "" {
		if !strings.Contains(*imageOverride, ":") && !strings.Contains(*imageOverride, "/") {
			// only the tag was provided
			*imageOverride = fmt.Sprintf("%s/%s/%s:%s", defaultRegistry, defaultOrganization, defaultImage, *imageOverride)
		} else if !containsRegistry(*imageOverride) {
			*imageOverride = fmt.Sprintf("%s/%s", defaultRegistry, *imageOverride)
		}

		if olm {
			image, err := replaceTagsWithDigests(imageOverride)
			if err != nil {
				digestErrors[image] = err
			}
		}
	}
}

func containsRegistry(image string) bool {
	re := regexp.MustCompile(`^[\w-\.]+\.[\w-\.]+\/[\w-\/@\.:]+$`)
	return re.MatchString(image)
}

// replaceTagsWithDigests replaces image tags with digests. This is needed for OLM installations.
// In case of error returns the provided image and the error
func replaceTagsWithDigests(imageOverride *string) (string, error) {
	if !strings.Contains(*imageOverride, "@sha256:") {
		prefix := strings.Split(*imageOverride, ":")[0]

		providedValue := *imageOverride
		digest, err := CraneDigest(providedValue)
		if err != nil {
			// set placeholder
			*imageOverride = prefix + "@" + missingDigestPlaceholder
			return providedValue, err
		}

		*imageOverride = prefix + "@" + digest
	}
	return *imageOverride, nil
}

func checkDigestErrors(digestErrors map[string]error) bool {
	if len(digestErrors) > 0 {
		fmt.Fprintln(os.Stderr, "WARNING: unable to retrieve the digest for some images. Please replace the placeholders in the YAML file.")
		fmt.Fprintln(os.Stderr, "Specific errors are:")
		for image, err := range digestErrors {
			fmt.Fprintf(os.Stderr, "- %s: %s\n", image, err.Error())
		}
		return true
	}
	return false
}

// in case of inconsistencies between the provided images and the selected installation type the user is warned
func checkInstallationTypeImagesMismatch(image, expectedRepo string, imageSetType imagesettype.ImageSetType) {

	// the check is performed only when using official Entando images
	if strings.HasPrefix(image, defaultRegistry+"/"+defaultOrganization+"/") {
		re := regexp.MustCompile(`^.+/(.+):.+$`)
		matches := re.FindStringSubmatch(image)
		if len(matches) == 2 {
			providedRepo := matches[1]
			if providedRepo != expectedRepo {
				fmt.Fprintf(os.Stderr, "WARNING: image-set-type is set to %s but the repository %s was provided. Expected repository should be %s\n", imageSetType, providedRepo, expectedRepo)
			}
		}
	}
}
