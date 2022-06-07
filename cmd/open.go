package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
	pkg "io.twasyl/devcore/pkg/utils"
)

func init() {
	rootCmd.AddCommand(buildOpenCommand())
}

func buildOpenCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "open",
		Short: "Opens a set of resources",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS != "darwin" {
				return errors.New(fmt.Sprintf("Operating system '%s' not currently supported for this operation", runtime.GOOS))
			}
			resource := args[0]
			if resource != "projects" && resource != "projects-dir" && resource != "servers" && resource != "servers-dir" {
				return errors.New(fmt.Sprintf("'%s' is unknown", resource))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			resource := args[0]
			if resource == "projects" || resource == "projects-dir" {
				resource = config.Config.ProjectsDir
			} else if resource == "servers" || resource == "servers-dir" {
				resource = config.Config.ServersDir
			}

			if _, err := os.Stat(resource); os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("The resource to open does not exist: %s", resource))
			}

			_, err := pkg.ExecCommand("open", resource)
			return err
		},
	}

	return command
}
