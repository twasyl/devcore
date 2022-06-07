package cmd

import (
	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
)

func init() {
	rootCmd.AddCommand(buildConfigCommand())
}

func buildConfigCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c"},
		Short:   "devcore config management",
	}

	command.AddCommand(buildRestoreDefaultToolsVersions())

	return command
}

func buildRestoreDefaultToolsVersions() *cobra.Command {
	command := &cobra.Command{
		Use:     "restore-default-tools-versions",
		Aliases: []string{"r"},
		Short:   "Restore default tools versions stored in configuration by the ones shipped with devcore",
		RunE: func(cmd *cobra.Command, args []string) error {

			if config.Config.DefaultToolsVersion == nil {
				config.Config.DefaultToolsVersion = make(map[string]string)
			}

			for _, tool := range config.DefaultSupportedTools {
				config.Config.DefaultToolsVersion[tool.Name] = tool.DefaultVersion
			}

			return config.Save()
		},
	}

	return command
}
