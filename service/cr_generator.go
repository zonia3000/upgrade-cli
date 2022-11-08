package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"upgrade-cli/common"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	defaultResourceName       = "my-app"
	apiVersion                = "app.entando.org/v1alpha1"
	EntandoAppNameEnv         = "ENTANDO_CLI_APPNAME"
	EntandoIngressHostNameEnv = "ENTANDO_CLI_INGRESS_HOST_NAME"
)

// GenerateCustomResource writes the CR in YAML format to the specified file or to the stdout if the filename is an empty string
// If the needsFix boolean flag is set to true a comment is added to the output to inform the user that some placeholders
// need to be replaced. Moreover the YAML syntax is broken to avoid accidental applies before human intervention.
func GenerateCustomResource(fileName string, entandoAppV2 *v1alpha1.EntandoAppV2, needsFix bool) error {

	entandoAppV2.APIVersion = apiVersion
	entandoAppV2.Kind = common.EntandoAppResourceName
	entandoAppV2.Name = defaultResourceName

	// set entandoAppName
	entandoAppName := os.Getenv(EntandoAppNameEnv)
	if entandoAppName == "" {
		return fmt.Errorf("the environment variable %s must be set", EntandoAppNameEnv)
	}
	entandoAppV2.Spec.EntandoAppName = entandoAppName

	// set ingressHostName
	ingressHostName := os.Getenv(EntandoIngressHostNameEnv)
	if ingressHostName == "" {
		return fmt.Errorf("the environment variable %s must be set", EntandoIngressHostNameEnv)
	}
	entandoAppV2.Spec.IngressHostName = ingressHostName

	yamlPrinter := printers.YAMLPrinter{}

	// set writer to file or to stdout if file is not specified
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

	writer.Write([]byte("---\n"))

	// write data to a buffer to be able to modify the result before writing it to the writer
	var buffer bytes.Buffer
	err := yamlPrinter.PrintObj(entandoAppV2, &buffer)

	if err != nil {
		return fmt.Errorf("unable to generate EntandoAppV2 manifest. %s", err.Error())
	}

	if needsFix {
		writer.Write(breakSyntax(buffer.Bytes()))
	} else {
		writer.Write(buffer.Bytes())
	}

	return nil
}

// breakSyntax removes the quotes around the error placeholders to break YAML syntax, in order to prevent accidental editing
func breakSyntax(bytes []byte) []byte {
	stringValue := string(bytes)
	errorRegexp := regexp.MustCompile(`("|')(ERROR: .*)("|')`)
	return []byte(errorRegexp.ReplaceAllString(stringValue, "$2 # FIXME"))
}
