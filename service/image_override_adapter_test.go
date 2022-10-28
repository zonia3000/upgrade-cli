package service

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/crane"
)

func TestAdaptImagesOverride(t *testing.T) {

	origDigest := craneDigest
	defer func() { craneDigest = origDigest }()

	craneDigest = func(ref string, opt ...crane.Option) (string, error) {
		return "sha256:94af0fb4525", nil
	}

	imagesOverride := make(map[string]string)

	imagesOverride["app-builder"] = "entando/app-builder:7.1.1-ENG-4277-PR-1413"
	imagesOverride["de-app"] = "registry.hub.docker.com/entando/entando-de-app-eap:7.1.1-ENGPM-493-PR-440"
	imagesOverride["keycloak"] = "entando/entando-keycloak@sha256:d550b07f5dd6"

	AdaptImagesOverride(imagesOverride)

	expectedAppBuilder := "registry.hub.docker.com/entando/app-builder@sha256:94af0fb4525"
	expectedDeApp := "registry.hub.docker.com/entando/entando-de-app-eap@sha256:94af0fb4525"
	expectedKeycloak := "registry.hub.docker.com/entando/entando-keycloak@sha256:d550b07f5dd6"

	if appBuilder := imagesOverride["app-builder"]; appBuilder != expectedAppBuilder {
		t.Fatalf("expected %s, found %s", expectedAppBuilder, appBuilder)
	}
	if deApp := imagesOverride["de-app"]; deApp != expectedDeApp {
		t.Fatalf("expected %s, found to %s", expectedDeApp, deApp)
	}
	if keycloak := imagesOverride["keycloak"]; keycloak != expectedKeycloak {
		t.Fatalf("expected %s, found %s", expectedKeycloak, keycloak)
	}
}
