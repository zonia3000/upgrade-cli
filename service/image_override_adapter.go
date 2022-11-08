package service

import (
	"fmt"
	"os"
	"strings"
	imagesettype "upgrade-cli/flag/image_set_type"
	"upgrade-cli/util/images"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/google/go-containerregistry/pkg/crane"
)

var CraneDigest = crane.Digest

const (
	missingDigestPlaceholder = "ERROR: <unable to fetch digest of: %s>"
)

// AdaptImagesOverride converts the format of the images provided by the user to full URL format
// Returns a bool that is true in case of errors in digests retrieval.
func AdaptImagesOverride(entandoAppV2 *v1alpha1.EntandoAppV2, olm bool) bool {

	imageSetType := imagesettype.ImageSetType(entandoAppV2.Spec.ImageSetType)

	digestErrors := make(map[string]error)

	for _, imageInfo := range images.EntandoImages {
		adaptImageOverride(entandoAppV2, imageInfo, imageSetType, olm, digestErrors)
	}

	return checkDigestErrors(digestErrors)
}

func adaptImageOverride(entandoAppV2 *v1alpha1.EntandoAppV2, imageInfo images.EntandoImageInfo, imageSetType imagesettype.ImageSetType, olm bool, digestErrors map[string]error) {
	imageOverride := imageInfo.GetImageOverride(entandoAppV2)

	if imageOverride != nil && *imageOverride != "" {
		if !strings.Contains(*imageOverride, ":") && !strings.Contains(*imageOverride, "/") {
			// only the tag was provided
			defaultImage := imageInfo.GetDefaultImage(imageSetType)
			*imageOverride = fmt.Sprintf("%s:%s", defaultImage, *imageOverride)
		} else if !images.ContainsRegistry(*imageOverride) {
			*imageOverride = fmt.Sprintf("%s/%s", images.DefaultRegistry, *imageOverride)
		}

		checkImageSetTypeMismatch(*imageOverride, imageInfo, imageSetType)

		if olm {
			err := replaceTagsWithDigests(imageOverride)
			if err != nil {
				digestErrors[imageInfo.ImageOverrideFlag] = err
			}
		}
	}
}

// replaceTagsWithDigests replaces image tags with digests. This is needed for OLM installations.
func replaceTagsWithDigests(imageOverride *string) error {
	if !strings.Contains(*imageOverride, "@sha256:") {
		prefix := strings.Split(*imageOverride, ":")[0]

		providedValue := *imageOverride
		digest, err := CraneDigest(providedValue)
		if err != nil {
			// set placeholder
			*imageOverride = fmt.Sprintf(missingDigestPlaceholder, providedValue)
			return err
		}

		*imageOverride = prefix + "@" + digest
	}
	return nil
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
func checkImageSetTypeMismatch(image string, imageInfo images.EntandoImageInfo, imageSetType imagesettype.ImageSetType) {

	// the check is performed only when using official Entando images
	if images.IsOfficialImage(image) && imageInfo.IsMultiImage {
		providedRepo := images.ExtractRepo(image)
		if providedRepo != "" {
			expectedRepo := images.ExtractRepo(imageInfo.GetDefaultImage(imageSetType))
			if providedRepo != expectedRepo {
				fmt.Fprintf(os.Stderr, "WARNING: image-set-type is set to %s but the repository %s was provided. Expected repository should be %s\n", imageSetType, providedRepo, expectedRepo)
			}
		}
	}
}
