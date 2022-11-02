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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prevent showing usage message when error happens in RunE func
		cmd.SilenceUsage = true

		file, _ := cmd.Flags().GetString(fileFlag)
		force, _ := cmd.Flags().GetBool(forceFlag)

		err := service.CreateEntandoApp(file, force)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Changes applied\n")

		var bar *progressbar.ProgressBar

		for {
			entandoApp, err := service.GetEntandoApp()
			if err != nil {
				return err
			}

			status, err := parseStatus(entandoApp)

			if err != nil {
				bar.Close()
				return err
			}

			if bar == nil {
				bar = progressbar.Default(int64(status.Total), "Upgrade in progress...")
				progressbar.OptionSetPredictTime(false)
			}
			bar.Set(status.Progress)

			if status.Progress == status.Total {
				bar.Close()
				fmt.Fprintf(os.Stderr, "Upgrade successfully completed")
				return nil
			}

			time.Sleep(1 * time.Second)
		}
	},
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
