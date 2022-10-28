package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
)

const defaultRegistry = "registry.hub.docker.com"

var craneDigest = crane.Digest

func AdaptImagesOverride(imagesOverride map[string]string) error {
	addMissingRegistry(imagesOverride)
	return replaceTagsWithDigests(imagesOverride)
}

func addMissingRegistry(imagesOverride map[string]string) {

	for key, image := range imagesOverride {
		re := regexp.MustCompile(`^[\w-\.]+\.[\w-\.]+\/[\w-\/@\.:]+$`)
		if !re.MatchString(image) {
			imagesOverride[key] = defaultRegistry + "/" + image
		}
	}
}

func replaceTagsWithDigests(imagesOverride map[string]string) error {

	for key, image := range imagesOverride {
		if !strings.Contains(image, "@sha256:") {
			digest, err := craneDigest(image)

			if err != nil {
				return fmt.Errorf("unable to retrieve digest for image %s.\n%s", image, err.Error())
			}

			imagesOverride[key] = strings.Split(image, ":")[0] + "@" + digest
		}
	}

	return nil
}
