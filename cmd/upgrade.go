package cmd

import (
	"fmt"
	"os"
	"time"
	"upgrade-cli/service"

	"github.com/entgigi/upgrade-operator.git/api/v1alpha1"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	fileFlag  = "file"
	forceFlag = "force"

	Succeeded = "Succeeded"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Apply EntandoAppV2 CR file",
	PreRun: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString(fileFlag)
		if file == "" {
			generateCRCmd.PreRun(cmd, args)
		}
		// If file flag is set, flags related to generation should not be set
		for componentFlag := range componentFlags {
			cmd.MarkFlagsMutuallyExclusive(fileFlag, componentFlag)
		}
		cmd.MarkFlagsMutuallyExclusive(fileFlag, versionFlag)
		cmd.MarkFlagsMutuallyExclusive(fileFlag, latestFlag)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prevent showing usage message when error happens in RunE func
		cmd.SilenceUsage = true

		fileName, _ := cmd.Flags().GetString(fileFlag)
		force, _ := cmd.Flags().GetBool(forceFlag)

		if fileName == "" {
			file, err := os.CreateTemp("", "entandoapp-cr")
			if err != nil {
				return err
			}

			fileName = file.Name()
			defer os.Remove(fileName)

			entandoApp, err := ParseEntandoAppFromCmd(cmd)
			if err != nil {
				return err
			}

			err = service.GenerateCustomResource(fileName, entandoApp)
			if err != nil {
				return err
			}
		}

		err := service.CreateEntandoApp(fileName, force)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Changes applied\n")

		return displayProgress()
	},
}

func displayProgress() error {
	var bar *progressbar.ProgressBar

	for {
		entandoApp, err := service.GetEntandoApp()
		if err != nil {
			if bar != nil {
				bar.Close()
			}
			return err
		}

		status, err := parseStatus(entandoApp)

		if err != nil {
			if bar != nil {
				bar.Close()
			}
			return err
		}

		if bar == nil {
			bar = newProgressbar(status.Total)
		}
		bar.Set(status.Progress)

		if status.Progress == status.Total {
			bar.Close()
			fmt.Fprintf(os.Stderr, "Upgrade successfully completed\n")
			return nil
		}

		time.Sleep(1 * time.Second)
	}
}

func newProgressbar(total int) *progressbar.ProgressBar {
	return progressbar.NewOptions(total,
		progressbar.OptionSetDescription("Upgrade in progress..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionShowCount(),
		progressbar.OptionFullWidth(),
	)
}

func parseStatus(entandoApp *v1alpha1.EntandoAppV2) (*v1alpha1.EntandoAppV2Status, error) {

	for _, condition := range entandoApp.Status.Conditions {
		if condition.Type == Succeeded && condition.Status == metav1.ConditionFalse {
			return nil, fmt.Errorf(condition.Message)
		}
	}

	return &entandoApp.Status, nil
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().Bool(forceFlag, false, "if set, the changes to the CR are applied even if the resource already exists")
	upgradeCmd.Flags().StringP(fileFlag, "f", "", "path to CR file")
}
