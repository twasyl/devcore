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
	rootCmd.AddCommand(buildDockerComposeCommand())
}

func buildDockerComposeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "docker-compose",
		Aliases: []string{"dc"},
		Short:   "docker compose related utilities",
	}

	command.AddCommand(buildDockerComposeContextCommand())

	return command
}

func buildDockerComposeContextCommand() *cobra.Command {
	context := config.DockerComposeContext{}

	command := &cobra.Command{
		Use:   "context",
		Short: "Manage docker compose context",
	}

	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Creates a docker compose context in the CLI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := config.Config.DockerCompose.FindContextByName(context.Name); config.IsDockerComposeContextNotFound(err) {
				if _, err := os.Stat(context.File); os.IsNotExist(err) {
					return errors.New(fmt.Sprintf("The file %s does not exist", context.File))
				}
			} else {
				return errors.New(fmt.Sprintf("A Docker compose context named '%s' already exists", context.Name))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Config.DockerCompose.AddContext(context)
			config.Config.DockerCompose.CurrentContext = context.Name
			return config.Save()
		},
	}

	createCommand.Flags().StringVarP(&context.Name, "name", "n", "", "The name of the context")
	createCommand.Flags().StringVarP(&context.Description, "description", "d", "", "The description of the context")
	createCommand.Flags().StringVarP(&context.File, "file", "f", "", "The docker compose file of the context")
	createCommand.MarkFlagRequired("name")
	createCommand.MarkFlagRequired("file")
	command.AddCommand(createCommand)

	var verbose bool
	listCommand := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists existing docker compose contexts",
		Run: func(cmd *cobra.Command, args []string) {
			for index, context := range config.Config.DockerCompose.Contexts {
				if verbose {
					if index == 0 {
						fmt.Println("")
					}
					if context.Name == config.Config.DockerCompose.CurrentContext {
						fmt.Println(fmt.Sprintf("* Name: %s", context.Name))
					} else {
						fmt.Println(fmt.Sprintf("  Name: %s", context.Name))
					}
					fmt.Println(fmt.Sprintf("  Description: %s", context.Description))
					fmt.Println(fmt.Sprintf("  File: %s", context.File))

					if index < len(config.Config.DockerCompose.Contexts)-1 {
						fmt.Println("")
					}

					if index == len(config.Config.DockerCompose.Contexts)-1 {
						fmt.Println("")
					}
				} else {
					if context.Name == config.Config.DockerCompose.CurrentContext {
						fmt.Println(fmt.Sprintf("* %s", context.Name))
					} else {
						fmt.Println(fmt.Sprintf("  %s", context.Name))
					}
				}
			}
		},
	}
	listCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display verbose contexts")
	command.AddCommand(listCommand)

	deleteCommand := &cobra.Command{
		Use:   "delete",
		Short: "Delete a docker compose context",
		Long:  "Delete a docker compose context from devcore without deleting the actual files",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToDelete := ""
			if len(args) == 1 {
				contextToDelete = args[0]
			} else if config.Config.DockerCompose.CurrentContext != "" {
				contextToDelete = config.Config.DockerCompose.CurrentContext
			} else {
				return errors.New("No docker compose context specified, neither a current one is set")
			}

			if c, err := config.Config.DockerCompose.FindContextByName(contextToDelete); config.IsDockerComposeContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Config.DockerCompose.DeleteContext(context)
			if context.Name == config.Config.DockerCompose.CurrentContext {
				config.Config.DockerCompose.CurrentContext = ""
			}
			err := config.Save()
			if err == nil {
				fmt.Println(fmt.Sprintf("Docker compose context '%s' deleted", context.Name))
			}
			return err
		},
	}
	command.AddCommand(deleteCommand)

	setCurrentCommand := &cobra.Command{
		Use:   "set-current",
		Short: "Set the current Docker compose context",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if c, err := config.Config.DockerCompose.FindContextByName(args[0]); config.IsDockerComposeContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Config.DockerCompose.CurrentContext = context.Name
			err := config.Save()
			if err == nil {
				fmt.Println(fmt.Sprintf("Current docker compose context set to '%s'", context.Name))
			}
			return err
		},
	}
	command.AddCommand(setCurrentCommand)

	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Starts a docker compose context",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToStart := ""
			if len(args) == 1 {
				contextToStart = args[0]
			} else if config.Config.DockerCompose.CurrentContext != "" {
				contextToStart = config.Config.DockerCompose.CurrentContext
			} else {
				return errors.New("No docker compose context specified, neither a current one is set")
			}

			if c, err := config.Config.DockerCompose.FindContextByName(contextToStart); config.IsDockerComposeContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("Starting context '%s'", context.Name))
			out, err := pkg.ExecCommand("docker", "compose", "-f", context.File, "up", "-d")
			if verbose {
				fmt.Println(out)
			}
			return err
		},
	}
	startCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	command.AddCommand(startCommand)

	stopCommand := &cobra.Command{
		Use:   "stop",
		Short: "Stops a docker compose context",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToStop := ""
			if len(args) == 1 {
				contextToStop = args[0]
			} else if config.Config.DockerCompose.CurrentContext != "" {
				contextToStop = config.Config.DockerCompose.CurrentContext
			} else {
				return errors.New("No docker compose context specified, neither a current one is set")
			}

			if c, err := config.Config.DockerCompose.FindContextByName(contextToStop); config.IsDockerComposeContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("Stopping context '%s'", context.Name))
			out, err := pkg.ExecCommand("docker", "compose", "-f", context.File, "down", "-v")
			if verbose {
				fmt.Println(out)
			}
			return err
		},
	}
	stopCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	command.AddCommand(stopCommand)

	openFolderCmd := &cobra.Command{
		Use:   "open-folder",
		Short: "Open the file explorer at the docker compose context",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToStop := ""
			if len(args) == 1 {
				contextToStop = args[0]
			} else if config.Config.DockerCompose.CurrentContext != "" {
				contextToStop = config.Config.DockerCompose.CurrentContext
			} else {
				return errors.New("No docker compose context specified, neither a current one is set")
			}

			if c, err := config.Config.DockerCompose.FindContextByName(contextToStop); config.IsDockerComposeContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if runtime.GOOS == "darwin" {
				_, err := pkg.ExecCommand("open", context.Dir())
				return err
			} else {
				return errors.New(fmt.Sprintf("%s not supported for this action", runtime.GOOS))
			}
		},
	}
	command.AddCommand(openFolderCmd)

	return command
}
