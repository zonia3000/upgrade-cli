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

func GenerateCustomResource(fileName string, entandoAppV2 *v1alpha1.EntandoAppV2) error {

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

	err := yamlPrinter.PrintObj(entandoAppV2, writer)
	if err != nil {
		return fmt.Errorf("unable to generate EntandoAppV2 manifest. %s", err.Error())
	}

	return nil
}
