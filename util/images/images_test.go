package images

import "testing"

func TestExtractRepoOLM(t *testing.T) {

	extractedRepo := ExtractRepo("registry.hub.docker.com/entando/entando-de-app-eap@sha256:94af0fb4525")
	expectedRepo := "entando-de-app-eap"
	if extractedRepo != expectedRepo {
		t.Fatalf("expected %s, found %s", expectedRepo, extractedRepo)
	}
}

func TestExtractRepoPlain(t *testing.T) {

	extractedRepo := ExtractRepo("registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.0")
	expectedRepo := "entando-de-app-wildfly"
	if extractedRepo != expectedRepo {
		t.Fatalf("expected %s, found %s", expectedRepo, extractedRepo)
	}
}

func TestExtractRepoNoVersion(t *testing.T) {

	extractedRepo := ExtractRepo("registry.hub.docker.com/entando/entando-de-app-wildfly")
	expectedRepo := "entando-de-app-wildfly"
	if extractedRepo != expectedRepo {
		t.Fatalf("expected %s, found %s", expectedRepo, extractedRepo)
	}
}

func TestContainsRegistry(t *testing.T) {

	imageFullUrl := "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.0"
	if !ContainsRegistry(imageFullUrl) {
		t.Fatalf("ContainsRegistry returned false for %s", imageFullUrl)
	}

	imageWithoutRegistry := "entando/entando-de-app-wildfly:7.1.0"
	if ContainsRegistry(imageWithoutRegistry) {
		t.Fatalf("ContainsRegistry returned true for %s", imageWithoutRegistry)
	}
}

func TestIsValidImageOverride(t *testing.T) {
	// Valid
	checkIsValidImageOverride(t, "7.1.0", true)
	checkIsValidImageOverride(t, "entando/entando-de-app-wildfly:7.1.0", true)
	checkIsValidImageOverride(t, "registry.hub.docker.com/entando/entando-de-app-wildfly:7.1.0", true)
	checkIsValidImageOverride(t, "registry.hub.docker.com/entando/entando-de-app-eap@sha256:94af0fb4525", true)
	// Invalid
	checkIsValidImageOverride(t, "entando-de-app-wildfly:7.1.0", false)
	checkIsValidImageOverride(t, "https://registry.hub.docker.com/r/entando/entando-de-app-wildfly", false)
}

func checkIsValidImageOverride(t *testing.T, imageOverride string, expected bool) {
	if IsValidImageOverride(imageOverride) != expected {
		t.Fatalf("IsValidImageOverride returned %v for %s", !expected, imageOverride)
	}
}
