package generate

import (
	"errors"
	"os"
	"strings"
	"testing"
	"upgrade-cli/service"

	"github.com/google/go-containerregistry/pkg/crane"
)

func TestGenerateSimpleCR(t *testing.T) {

	os.Setenv(service.EntandoAppNameEnv, "my-entando-app")
	os.Setenv(service.EntandoIngressHostNameEnv, "quickstart.10.11.91.88.nip.io")

	testFile, _ := os.CreateTemp("", "generate-cr-test")
	defer os.Remove(testFile.Name())

	GenerateCRCmd.SetArgs([]string{"generate", "-o", testFile.Name(), "-v", "7.1.0", "--operator-mode", "Plain"})

	err := GenerateCRCmd.Execute()

	if err != nil {
		t.Fatalf(err.Error())
	}

	bytes, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatalf(err.Error())
	}

	fileContent := string(bytes)

	assertYamlField(t, fileContent, "apiVersion", "app.entando.org/v1alpha1")
	assertYamlField(t, fileContent, "kind", "EntandoAppV2")
	assertYamlField(t, fileContent, "version", "7.1.0")
	assertYamlField(t, fileContent, "entandoAppName", "my-entando-app")
	assertYamlField(t, fileContent, "ingressHostName", "quickstart.10.11.91.88.nip.io")
}

func TestGenerateOlmCRWithPlaceholders(t *testing.T) {

	os.Setenv(service.EntandoAppNameEnv, "my-entando-app")
	os.Setenv(service.EntandoIngressHostNameEnv, "quickstart.10.11.91.88.nip.io")

	testFile, _ := os.CreateTemp("", "generate-cr-test")
	defer os.Remove(testFile.Name())

	origDigest := service.CraneDigest
	defer func() { service.CraneDigest = origDigest }()

	service.CraneDigest = func(ref string, opt ...crane.Option) (string, error) {
		if strings.HasSuffix(ref, "invalid-tag") {
			return "", errors.New("manifest unknown")
		} else {
			return "sha256:94af0fb4525", nil
		}
	}

	GenerateCRCmd.SetArgs([]string{"generate", "-o", testFile.Name(), "-v", "v7.1.0", "--operator-mode", "OLM",
		"--image-de-app", "7.1.0-fix1", "--image-app-builder", "invalid-tag"})

	err := GenerateCRCmd.Execute()

	if err != nil {
		t.Fatalf(err.Error())
	}

	bytes, err := os.ReadFile(testFile.Name())
	if err != nil {
		t.Fatalf(err.Error())
	}

	fileContent := string(bytes)

	wrongImageResult := "imageOverride: ERROR: <unable to fetch digest of: registry.hub.docker.com/entando/app-builder:invalid-tag> # FIXME"
	if !strings.Contains(fileContent, wrongImageResult) {
		t.Fatalf("Generated doesn't contain placeholders warning")
	}

	assertYamlField(t, fileContent, "imageOverride", "registry.hub.docker.com/entando/entando-de-app-eap@sha256:94af0fb4525")
}

func assertYamlField(t *testing.T, fileContent, key, expectedValue string) {
	if !strings.Contains(fileContent, key+": "+expectedValue) {
		t.Fatalf("Generated CR doesn't contain expected %s\n%s", key, fileContent)
	}
}

func TestInvalidImageOverrideFormat(t *testing.T) {
	GenerateCRCmd.SetArgs([]string{"generate", "-v", "v7.1.0", "--image-de-app", "foo:bar:foo"})

	err := GenerateCRCmd.Execute()

	expectedErrorMessage := "invalid format for image override flag 'foo:bar:foo'. It should be <image>:<tag> or <tag>"

	if err == nil {
		t.Fatalf("an error was expected")
	} else if expectedErrorMessage != err.Error() {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
