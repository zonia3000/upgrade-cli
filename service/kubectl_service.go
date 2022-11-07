package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	operatormode "upgrade-cli/flag/operator_mode"
	"upgrade-cli/util/sys/spawn"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	kubectlBaseCommandEnv  = "ENTANDO_CLI_KUBECTL_COMMAND"
	entandoAppResourceName = "EntandoAppV2"
	operatorDeploymentType = "ENTANDO_K8S_OPERATOR_DEPLOYMENT_TYPE"
)

// CreateEntandoApp sends the CR creation request to the cluster.
// If force is set to true an existing resource will be overwritten
func CreateEntandoApp(fileName string, force bool) error {

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %s doesn't exist", fileName)
	}

	baseCmd, args, err := getKubectlBaseCommand()
	if err != nil {
		return err
	}

	var kubectlCmd string
	if force {
		kubectlCmd = "apply"
	} else {
		kubectlCmd = "create"
	}

	args = append(args, kubectlCmd, "-f", fileName)

	output, err := spawn.Spawn(nil,
		*baseCmd,
		args,
		spawn.Environ{},
		spawn.Options{
			WithSudo:      false,
			CaptureStdout: true,
			CaptureStderr: true,
		},
	)

	if err != nil && len(output.Stderr) > 0 {
		if !force && strings.Contains(output.Stderr, "AlreadyExists") {
			return fmt.Errorf("resource already exists. You can overwrite it using the --force flag")
		}
		return fmt.Errorf("error creating the resource: %s", output.Stderr)
	}

	return err
}

// GetEntandoApp retrieves the EntandoAppV2 resource from the cluster
func GetEntandoApp() (*v1alpha1.EntandoAppV2, error) {
	baseCmd, args, err := getKubectlBaseCommand()
	if err != nil {
		return nil, err
	}

	args = append(args, "get", entandoAppResourceName, "-o", "yaml")

	output, err := spawn.Spawn(nil,
		*baseCmd,
		args,
		spawn.Environ{},
		spawn.Options{
			WithSudo:      false,
			CaptureStdout: true,
		},
	)
	if err != nil {
		return nil, err
	}

	return parseEntandoAppV2(output.Stdout)
}

func parseEntandoAppV2(stdout string) (*v1alpha1.EntandoAppV2, error) {

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	// Decode the YAML to an object
	entandoApps := v1alpha1.EntandoAppV2List{}
	_, _, err := s.Decode([]byte(stdout), nil, &entandoApps)
	if err != nil {
		return nil, err
	}

	if len(entandoApps.Items) == 0 {
		return nil, fmt.Errorf("resource of type %s not found", entandoAppResourceName)
	}
	if len(entandoApps.Items) > 1 {
		return nil, fmt.Errorf("found multiple resources of type %s", entandoAppResourceName)
	}

	return &entandoApps.Items[0], nil
}

// GetOperatorMode retrieves the OperatorMode from the cluster
// It reads the related environment variable inside entando-operator deployment spec
func GetOperatorMode() (operatormode.OperatorMode, error) {
	baseCmd, args, err := getKubectlBaseCommand()
	if err != nil {
		return operatormode.Auto, err
	}

	jsonPath := "jsonpath='{.spec.template.spec.containers[0].env[?(@.name == \"" + operatorDeploymentType + "\")].value}'"
	args = append(args, "get", "deploy", "entando-operator", "-o", jsonPath)

	output, err := spawn.Spawn(nil,
		*baseCmd,
		args,
		spawn.Environ{},
		spawn.Options{
			WithSudo:      false,
			CaptureStdout: true,
		},
	)
	if err != nil {
		return operatormode.Auto, fmt.Errorf("unable to retrieve the operator mode from the deployment: %s", err.Error())
	}

	parsedMode := strings.Trim(output.Stdout, "'")

	switch parsedMode {
	case "olm":
		return operatormode.OLM, nil
	case "helm":
		return operatormode.Plain, nil
	default:
		return operatormode.Auto, fmt.Errorf("unable to retrieve the operator mode from the deployment.\nUnexpected value for %s: %s", operatorDeploymentType, parsedMode)
	}
}

// getKubectlBaseCommand returns the base kubectl command parsed from the related environment variable
// and converted in the format required by the spawn.Spawn function
func getKubectlBaseCommand() (*string, []interface{}, error) {

	kubectlBaseCmd := os.Getenv(kubectlBaseCommandEnv)
	if kubectlBaseCmd == "" {
		return nil, nil, fmt.Errorf("the environment variable %s must be set", kubectlBaseCommandEnv)
	}

	parts := strings.Split(kubectlBaseCmd, " ")

	var kubectlArgs []interface{}
	for i := 1; i < len(parts); i = i + 1 {
		kubectlArgs = append(kubectlArgs, parts[i])
	}

	return &parts[0], kubectlArgs, nil
}
