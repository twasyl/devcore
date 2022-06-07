package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
)

func init() {
	rootCmd.AddCommand(toolsCommand)

	validArgs := make([]string, 0)
	for _, tool := range config.DefaultSupportedTools {
		validArgs = append(validArgs, tool.Name)
	}
	toolsInstallCmd.ValidArgs = validArgs

	toolsCommand.AddCommand(toolsInstallCmd)
	toolsInstallCmd.PersistentFlags().StringVarP(&toolVersion, "version", "v", "", "The version of the tool to install")

	toolsCommand.AddCommand(toolsListCmd)
	toolsListCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
}

var tool config.Tool
var toolVersion string
var verbose bool

var toolsCommand = &cobra.Command{
	Use:   "tools",
	Short: "Provide facilitators for developers' tools",
}

var toolsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a tool",
	Args:  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		tool, err = config.FindTool(args[0])
		if config.IsToolNotFound(err) {
			return err
		}

		if toolVersion == "" {
			toolVersion = tool.DefaultVersion
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(fmt.Sprintf("Installing %s %s", tool.Name, toolVersion))
		return tool.Install(toolVersion)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("%s installed successfully", tool.Name))
	},
}

var toolsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List supported tools",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Supported tools:")
		for _, tool := range config.DefaultSupportedTools {
			fmt.Printf("  - %s\n", tool.Name)
			if verbose {
				fmt.Printf("    Default version: %s\n", tool.DefaultVersion)
				fmt.Printf("    Description: %s\n", tool.Description)
			}
		}
	},
}
