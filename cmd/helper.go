package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var helperCmd = &cobra.Command{
		Use:   "helper",
		Short: "Utilities to help you manage tools & env",
	}

	rootCmd.AddCommand(helperCmd)

	deletionConfirmed := false

	var cleanDockerCmd = &cobra.Command{
		Use:   "uninstall-docker",
		Short: "Completely remove docker from your system",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print("Are you sure you want to delete Docker from your system? [yN] ")
			prompt := bufio.NewReader(os.Stdin)
			answer, err := prompt.ReadString('\n')
			if err != nil {
				return err
			}

			answer = strings.TrimSpace(answer)
			deletionConfirmed = answer == "y" || answer == "Y"
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if deletionConfirmed {

			} else {
				fmt.Println("Docker deletion aborted.")
			}
			return nil
		},
	}

	helperCmd.AddCommand(cleanDockerCmd)
}
