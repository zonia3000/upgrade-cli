package images

import (
	"fmt"
	"regexp"
	"strings"
	imagesettype "upgrade-cli/flag/image_set_type"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
)

const (
	DefaultRegistry     = "registry.hub.docker.com"
	DefaultOrganization = "entando"

	defaultAppBuilderImage                 = "app-builder"
	defaultDeAppEapImage                   = "entando-de-app-eap"
	defaultDeAppWildflyImage               = "entando-de-app-wildfly"
	defaultComponentManagerImage           = "entando-component-manager"
	defaultKeycloakImage                   = "entando-keycloak"
	defaultRedHatSsoImage                  = "entando-redhat-sso"
	defaultK8sServiceImage                 = "entando-k8s-service"
	defaultK8sPluginControllerImage        = "entando-k8s-plugin-controller"
	defaultK8sAppPluginLinkControllerImage = "entando-k8s-app-plugin-link-controller"
)

// EntandoImageInfo contains information used to parse image overrides and retrieve default images
type EntandoImageInfo struct {
	ComponentName string
	// name of the flag used to specify the image override
	ImageOverrideFlag string
	// true if the component supports different images depending on the imageSetType
	IsMultiImage bool
	// returns the default image according to the specified imageSetType
	// if the component is not multi-image the function returns a constant value
	GetDefaultImage func(imageSetType imagesettype.ImageSetType) string
	// returns the reference to the related ImageOverride field in EntandoAppV2
	GetImageOverride func(entandoApp *v1alpha1.EntandoAppV2) *string
}

// list of all the Entando components images
var EntandoImages = []EntandoImageInfo{

	mkEntandoComponentInfoMultiImage("DeApp", "image-de-app",
		defaultDeAppWildflyImage, defaultDeAppEapImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.DeApp.ImageOverride
		}),

	mkEntandoComponentInfoSingleImage("AppBuilder", "image-app-builder",
		defaultAppBuilderImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.AppBuilder.ImageOverride
		}),

	mkEntandoComponentInfoSingleImage("ComponentManager", "image-component-manager",
		defaultComponentManagerImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.ComponentManager.ImageOverride
		}),

	mkEntandoComponentInfoMultiImage("Keycloak", "image-keycloak",
		defaultKeycloakImage, defaultRedHatSsoImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.Keycloak.ImageOverride
		}),

	mkEntandoComponentInfoSingleImage("K8sService", "image-k8s-service",
		defaultK8sServiceImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.K8sService.ImageOverride
		}),

	mkEntandoComponentInfoSingleImage("K8sPluginController", "image-k8s-plugin-controller",
		defaultK8sPluginControllerImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.K8sPluginController.ImageOverride
		}),

	mkEntandoComponentInfoSingleImage("K8sAppPluginLinkController", "image-k8s-app-plugin-link-controller",
		defaultK8sAppPluginLinkControllerImage,
		func(entandoApp *v1alpha1.EntandoAppV2) *string {
			return &entandoApp.Spec.K8sAppPluginLinkController.ImageOverride
		}),
}

func mkEntandoComponentInfoSingleImage(name, flag, defaultRepo string, getImageOverride func(entandoApp *v1alpha1.EntandoAppV2) *string) EntandoImageInfo {
	return EntandoImageInfo{
		ComponentName:     name,
		ImageOverrideFlag: flag,
		IsMultiImage:      false,
		GetDefaultImage:   mkGetDefaultSingleImage(defaultRepo),
		GetImageOverride:  getImageOverride,
	}
}

func mkEntandoComponentInfoMultiImage(name, flag, defaultRepoCommunity, defaultRepoCertified string, getImageOverride func(entandoApp *v1alpha1.EntandoAppV2) *string) EntandoImageInfo {
	return EntandoImageInfo{
		ComponentName:     name,
		ImageOverrideFlag: flag,
		IsMultiImage:      true,
		GetDefaultImage:   mkGetDefaultMultiImage(defaultRepoCommunity, defaultRepoCertified),
		GetImageOverride:  getImageOverride,
	}
}

func mkGetDefaultMultiImage(repoCommunity, repoCertified string) func(imageSetType imagesettype.ImageSetType) string {
	return func(imageSetType imagesettype.ImageSetType) string {
		if imageSetType == imagesettype.Community {
			return mkDefaultImage(repoCommunity)
		}
		return mkDefaultImage(repoCertified)
	}
}

func mkGetDefaultSingleImage(repo string) func(imageSetType imagesettype.ImageSetType) string {
	return func(imageSetType imagesettype.ImageSetType) string {
		return mkDefaultImage(repo)
	}
}

func mkDefaultImage(repo string) string {
	return fmt.Sprintf("%s/%s/%s", DefaultRegistry, DefaultOrganization, repo)
}

// ExtractRepo extracts repository name from image full URL
func ExtractRepo(image string) string {
	re := regexp.MustCompile(`^.+/([^@:]+)(?:@sha256)?:?.*$`)
	matches := re.FindStringSubmatch(image)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// IsOfficialImage returns true if the provided image is an official Entando image
func IsOfficialImage(image string) bool {
	return strings.HasPrefix(image, DefaultRegistry+"/"+DefaultOrganization+"/")
}

// ContainsRegistry returns true if the provided image contains a registry
func ContainsRegistry(image string) bool {
	re := regexp.MustCompile(`^[\w-\.]+\.[\w-\.]+\/[\w-\/@\.:]+$`)
	return re.MatchString(image)
}

// IsValidImageOverride returns true if the provided value can be used as image override flag
// Accepted values are:
// - <tag>
// - <organization>/<repo>:<tag>
// - <organization>/<repo>@sha256:<sha>
// - <registry>/<organization>/<repo>:<tag>
// - <registry>/<organization>/<repo>@sha256:<sha>
func IsValidImageOverride(imageOverride string) bool {
	re := regexp.MustCompile(`^([\w-\.]+\/)?([\w-]+\/[\w-]+(@sha256)?)?:?[\w-\.]+$`)
	return re.MatchString(imageOverride)
}
