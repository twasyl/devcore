package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
)

func init() {
	rootCmd.AddCommand(serversCommand)

	validArgs := make([]string, 0)
	for _, server := range config.DefaultSupportedServers {
		validArgs = append(validArgs, server.Name)
	}
	serversInstallCmd.ValidArgs = validArgs

	serversCommand.AddCommand(serversInstallCmd)
	serversInstallCmd.PersistentFlags().StringVarP(&serverVersion, "version", "v", "", "The version of the server to install")
}

var server config.Server
var serverVersion string

var serversCommand = &cobra.Command{
	Use:   "servers",
	Short: "Provide facilitators for developers' servers",
}

var serversInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a server",
	Args:  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		server, err = config.FindServer(args[0])
		if config.IsServerNotFound(err) {
			return err
		}

		if serverVersion == "" {
			serverVersion = server.DefaultVersion
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(fmt.Sprintf("Installing %s %s", server.Name, serverVersion))
		return server.Install(serverVersion)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%s installed successfully", server.Name))
	},
}
