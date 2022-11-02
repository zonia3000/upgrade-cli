package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"upgrade-cli/util/sys/spawn"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	kubectlBaseCommandEnv  = "ENTANDO_KUBECTL_BASE_COMMAND"
	entandoAppResourceName = "EntandoAppV2"
)

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

	_, err = spawn.Spawn(nil,
		*baseCmd,
		args,
		spawn.Environ{},
		spawn.Options{
			WithSudo:      false,
			CaptureStdout: true,
		},
	)

	return err
}

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

	return parseResource(output.Stdout)
}

func parseResource(stdout string) (*v1alpha1.EntandoAppV2, error) {

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
