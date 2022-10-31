package service

import (
	"testing"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/google/go-containerregistry/pkg/crane"
)

func TestAdaptImagesOverride(t *testing.T) {

	origDigest := craneDigest
	defer func() { craneDigest = origDigest }()

	craneDigest = func(ref string, opt ...crane.Option) (string, error) {
		return "sha256:94af0fb4525", nil
	}

	entandoAppV2 := v1alpha1.EntandoAppV2{}

	entandoAppV2.Spec.AppBuilder.ImageOverride = "entando/app-builder:7.1.1-ENG-4277-PR-1413"
	entandoAppV2.Spec.DeApp.ImageOverride = "registry.hub.docker.com/entando/entando-de-app-eap:7.1.1-ENGPM-493-PR-440"
	entandoAppV2.Spec.Keycloak.ImageOverride = "entando/entando-keycloak@sha256:d550b07f5dd6"

	AdaptImagesOverride(&entandoAppV2)

	expectedAppBuilder := "registry.hub.docker.com/entando/app-builder@sha256:94af0fb4525"
	expectedDeApp := "registry.hub.docker.com/entando/entando-de-app-eap@sha256:94af0fb4525"
	expectedKeycloak := "registry.hub.docker.com/entando/entando-keycloak@sha256:d550b07f5dd6"

	if appBuilder := entandoAppV2.Spec.AppBuilder.ImageOverride; appBuilder != expectedAppBuilder {
		t.Fatalf("expected %s, found %s", expectedAppBuilder, appBuilder)
	}
	if deApp := entandoAppV2.Spec.DeApp.ImageOverride; deApp != expectedDeApp {
		t.Fatalf("expected %s, found to %s", expectedDeApp, deApp)
	}
	if keycloak := entandoAppV2.Spec.Keycloak.ImageOverride; keycloak != expectedKeycloak {
		t.Fatalf("expected %s, found %s", expectedKeycloak, keycloak)
	}
}
