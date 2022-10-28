package service

import (
	"fmt"
	"os"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	defaultAppName = "my-app"
	apiVersion     = "app.entando.org/v1alpha1"
	kind           = "EntandoAppV2"
)

func GenerateCustomResource(fileName string, version string, imagesOverride map[string]string) error {
	entandoAppV2 := v1alpha1.EntandoAppV2{}

	entandoAppV2.APIVersion = apiVersion
	entandoAppV2.Kind = kind
	entandoAppV2.Name = defaultAppName
	entandoAppV2.Spec.Version = version

	if len(imagesOverride) > 0 {
		entandoAppV2.Spec.ImagesOverride = imagesOverride
	}

	yamlPrinter := printers.YAMLPrinter{}

	crFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("unable to create file %s. %s", fileName, err.Error())
	}

	defer crFile.Close()

	err = yamlPrinter.PrintObj(&entandoAppV2, crFile)
	if err != nil {
		return fmt.Errorf("unable to generate EntandoAppV2 manifest. %s", err.Error())
	}

	return nil
}
