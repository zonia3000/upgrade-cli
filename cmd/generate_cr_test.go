package cmd

import (
	"os"
	"strings"
	"testing"
)

const testFile = "generate-cr-test.yml"

func TestSimpleCRGenerated(t *testing.T) {

	defer os.Remove(testFile)

	rootCmd.SetArgs([]string{"generate-cr", "-f", testFile, "-v", "v7.1.0"})

	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf(err.Error())
	}

	bytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fileContent := string(bytes)

	if !strings.Contains(fileContent, "apiVersion: app.entando.org/v1alpha1") {
		t.Fatalf("Generated CR doesn't contain expected apiVersion\n%s", fileContent)
	}
	if !strings.Contains(fileContent, "kind: EntandoAppV2") {
		t.Fatalf("Generated CR doesn't contain expected kind\n%s", fileContent)
	}
	if !strings.Contains(fileContent, "version: v7.1.0") {
		t.Fatalf("Generated CR doesn't contain expected version\n%s", fileContent)
	}
}

func TestInvalidImageOverrideFormat(t *testing.T) {
	rootCmd.SetArgs([]string{"generate-cr", "-f", testFile, "-v", "v7.1.0", "-i", "foo"})

	err := rootCmd.Execute()

	expectedErrorMessage := "invalid format for image override flag 'foo'. It should be <component-name>=<image>:<tag>"

	if err == nil {
		t.Fatalf("an error was expected")
	} else if expectedErrorMessage != err.Error() {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
