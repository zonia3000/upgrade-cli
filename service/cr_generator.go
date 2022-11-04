package service

import (
	"fmt"
	"io"
	"os"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	defaultResourceName = "my-app"
	apiVersion          = "app.entando.org/v1alpha1"
	kind                = "EntandoAppV2"
	EntandoAppNameEnv   = "ENTANDO_APPNAME"
)

// GenerateCustomResource writes the CR in YAML format to the specified file or to the stdout if the filename is an empty string
// If the needsFix boolean flag is set to true a comment is added to the output to inform the user that some placeholders
// need to be replaced. Moreover the YAML syntax is broken to avoid accidental applies before human intervention.
func GenerateCustomResource(fileName string, entandoAppV2 *v1alpha1.EntandoAppV2, needsFix bool) error {

	entandoAppV2.APIVersion = apiVersion
	entandoAppV2.Kind = kind
	entandoAppV2.Name = defaultResourceName

	entandoAppName := os.Getenv(EntandoAppNameEnv)
	if entandoAppName == "" {
		return fmt.Errorf("the environment variable %s must be set", EntandoAppNameEnv)
	}
	entandoAppV2.Spec.EntandoAppName = entandoAppName

	yamlPrinter := printers.YAMLPrinter{}

	var writer io.Writer
	if fileName == "" {
		writer = os.Stdout
	} else {
		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("unable to create file %s. %s", fileName, err.Error())
		}
		defer file.Close()
		writer = file
	}

	if needsFix {
		printFixme(writer)
	}

	err := yamlPrinter.PrintObj(entandoAppV2, writer)
	if err != nil {
		return fmt.Errorf("unable to generate EntandoAppV2 manifest. %s", err.Error())
	}

	return nil
}

func printFixme(writer io.Writer) {
	// adding the '<' char to break the syntax
	writer.Write([]byte("< ######################################## >\n"))
	writer.Write([]byte("< #                FIXME                 # >\n"))
	writer.Write([]byte("< #   Please replace the placeholders    # >\n"))
	writer.Write([]byte("< # Remember also to remove this comment # >\n"))
	writer.Write([]byte("< ######################################## >\n\n"))
}
