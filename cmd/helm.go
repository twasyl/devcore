package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	du "io.twasyl/devcore/pkg/utils"
)

func init() {
	rootCmd.AddCommand(helmCmd)

	helmCmd.AddCommand(helmInstallCmd)
	helmCmd.PersistentFlags().StringVarP(&helmNamespace, "namespace", "n", "", "The namespace to use when executing Helm commands")

	helmCmd.AddCommand(helmUninstallCmd)
}

var helmNamespace string
var helmReleaseName string
var helmChartName string

var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Provide Helm utilities",
	Long:  "helm command provide utilities related to helm",
}

var helmInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a Helm chart",
	Args:  cobra.ExactArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if helmNamespace != "" {
			output, err := du.ExecCommand("kubectl", "create", "namespace", helmNamespace)
			if err != nil {
				log.Fatal(output)
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		args = []string{"install", args[0], args[1]}
		if helmNamespace != "" {
			args = append(args, "-n", helmNamespace)
		}
		output, err := du.ExecCommand("helm", args...)
		fmt.Print(output)
		if err != nil {
			return err
		}
		return nil
	},
}

var helmUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a Helm release",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		args = []string{"uninstall", args[0]}
		if helmNamespace != "" {
			args = append(args, "-n", helmNamespace)
		}
		output, err := du.ExecCommand("helm", args...)
		fmt.Print(output)
		if err != nil {
			return err
		}
		return nil
	},
}
