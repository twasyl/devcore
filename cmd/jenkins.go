package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
	du "io.twasyl/devcore/pkg/utils"
	pkg "io.twasyl/devcore/pkg/utils"
)

func init() {
	rootCmd.AddCommand(buildJenkinsCommand())
}

func buildJenkinsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "jenkins",
		Short: "Interact with Jenkins",
	}

	command.AddCommand(buildJenkinsCliCommand())
	command.AddCommand(buildJenkinsContextCommand())
	return command
}

func buildJenkinsCliCommand() *cobra.Command {
	var jenkinsUrl string
	var jenkinsCreds string
	var useWebsockets bool

	var command = &cobra.Command{
		Use:   "cli",
		Short: "Relates to Jenkins CLI",
	}
	command.PersistentFlags().StringVarP(&jenkinsUrl, "url", "u", "http://127.0.0.1:8080", "URL of the Jenkins instance.")

	getCommand := &cobra.Command{
		Use:   "get",
		Short: "Download the Jenkins CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := du.DownloadFile(fmt.Sprintf("%s/jnlpJars/jenkins-cli.jar", jenkinsUrl), "jenkins-cli.jar")
			if err == nil {
				fmt.Println("The Jenkins CLI has been downloaded.")
				path, _ := filepath.Abs("jenkins-cli.jar")
				config.Config.Jenkins.Cli = path
				config.Save()
			}
			return err
		},
	}
	command.AddCommand(getCommand)

	execCommand := &cobra.Command{
		Use:   "exec",
		Short: "Execute Jenkins CLI commands",
		Long:  "Execute Jenkins CLI commands by automatically adding information connection to the call of jenkins-cli.jar",
		Args:  cobra.ArbitraryArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.Config.Jenkins.Cli == "" {
				return errors.New("No jenkins-cli.jar registered. Use 'jenkins cli get' first")
			}

			info, err := os.Stat(config.Config.Jenkins.Cli)
			if os.IsNotExist(err) {
				return err
			}
			if info.IsDir() {
				return errors.New("The jenkins-cli.jar is a directory")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs := []string{"-jar", "jenkins-cli.jar", "-s", jenkinsUrl}

			if useWebsockets {
				cliArgs = append(cliArgs, "-webSocket")
			}

			if jenkinsCreds != "" {
				cliArgs = append(cliArgs, "-auth", jenkinsCreds)
			}

			cliArgs = append(cliArgs, args...)
			output, err := du.ExecCommand("java", cliArgs...)
			fmt.Println(output)
			return err
		},
	}
	execCommand.Flags().StringVarP(&jenkinsCreds, "auth", "a", "", "Credentials used to connect to Jenkins")
	execCommand.Flags().BoolVarP(&useWebsockets, "websockets", "w", false, "Interact with Jenkins using websockets")
	command.AddCommand(execCommand)

	return command
}

func buildJenkinsContextCommand() *cobra.Command {
	context := config.JenkinsContext{}
	context.Pid = -1

	var command = &cobra.Command{
		Use:     "context",
		Short:   "Relates to Jenkins context",
		Aliases: []string{"ctx"},
	}
	command.PersistentFlags().StringVarP(&context.Name, "name", "n", "", "The name of the Jenkins context")

	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Creates a Jenkins context in the CLI",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := config.Config.Jenkins.FindContextByName(context.Name); config.IsJenkinsContextNotFound(err) {
				if _, err := os.Stat(context.War); os.IsNotExist(err) {
					return errors.New(fmt.Sprintf("The file %s does not exist", context.War))
				}
			} else {
				return errors.New(fmt.Sprintf("A Jenkins context named '%s' already exists", context.Name))
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Config.Jenkins.AddContext(context)
			config.Config.Jenkins.CurrentContext = context.Name
			return config.Save()
		},
	}

	createCommand.Flags().StringVarP(&context.Description, "description", "d", "", "Description of this context")
	createCommand.Flags().StringVarP(&context.War, "war", "w", "", "The Jenkins war to use")
	createCommand.Flags().StringVar(&context.JenkinsHome, "jenkins-home", "", "The folder to use as Jenkins home")
	createCommand.Flags().StringVar(&context.JavaHome, "java-home", "", "The Java home to use with this context. If unspecified, the default JAVA_HOME of the system will be used")
	createCommand.Flags().StringArrayVar(&context.Options, "option", nil, "The option to pass to Jenkins at startup. Use multiple times for multiple options")
	createCommand.Flags().StringArrayVar(&context.JVMOptions, "jvm-option", nil, "The JVM option to pass to Jenkins at startup. Use multiple times for multiple options")
	createCommand.MarkFlagRequired("name")
	createCommand.MarkFlagRequired("war")
	command.AddCommand(createCommand)

	var verbose bool
	listCommand := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists existing Jenkins contexts",
		Run: func(cmd *cobra.Command, args []string) {
			for index, context := range config.Config.Jenkins.Contexts {
				if verbose {
					if index == 0 {
						fmt.Println("")
					}
					if context.Name == config.Config.Jenkins.CurrentContext {
						fmt.Println(fmt.Sprintf("* Name: %s", context.Name))
					} else {
						fmt.Println(fmt.Sprintf("  Name: %s", context.Name))
					}
					fmt.Println(fmt.Sprintf("  Description: %s", context.Description))
					fmt.Println(fmt.Sprintf("  War: %s", context.War))
					fmt.Println(fmt.Sprintf("  Jenkins home: %s", context.JenkinsHome))
					fmt.Println(fmt.Sprintf("  Java home: %s", context.JavaHome))
					fmt.Println(fmt.Sprintf("  Options: %s", context.Options))
					fmt.Println(fmt.Sprintf("  PID: %d", context.Pid))

					if index < len(config.Config.Jenkins.Contexts)-1 {
						fmt.Println("")
						fmt.Println("")
					}
				} else {
					if context.Name == config.Config.Jenkins.CurrentContext {
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
		Short: "Delete a Jenkins context",
		Long:  "Delete a Jenkins context from devcore without deleting the actual files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, context := range args {
				if c, err := config.Config.Jenkins.FindContextByName(context); config.IsJenkinsContextNotFound(err) {
					fmt.Println(fmt.Sprintf("Context %s not found.\n", context))
				} else {
					err := config.Config.Jenkins.DeleteContext(c)
					if err != nil {
						return err
					}
					fmt.Println(fmt.Sprintf("Jenkins context '%s' deleted", c.Name))
					if c.Name == config.Config.Jenkins.CurrentContext {
						config.Config.Jenkins.CurrentContext = ""
					}
				}
			}

			err := config.Save()
			return err
		},
	}
	command.AddCommand(deleteCommand)

	setCurrentCommand := &cobra.Command{
		Use:   "set-current",
		Short: "Set the current Jenkins context",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if c, err := config.Config.Jenkins.FindContextByName(args[0]); config.IsJenkinsContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Config.Jenkins.CurrentContext = context.Name
			err := config.Save()
			if err == nil {
				fmt.Println(fmt.Sprintf("Current Jenkins context set to '%s'", context.Name))
			}
			return err
		},
	}
	command.AddCommand(setCurrentCommand)

	additionalJvmOptions := []string{}
	startCommand := &cobra.Command{
		Use:   "start",
		Short: "Starts a Jenkins context",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToStart := ""
			if len(args) == 1 {
				contextToStart = args[0]
			} else if config.Config.Jenkins.CurrentContext != "" {
				contextToStart = config.Config.Jenkins.CurrentContext
			} else {
				return errors.New("No Jenkins context specified, neither a current one is set")
			}

			if c, err := config.Config.Jenkins.FindContextByName(contextToStart); config.IsJenkinsContextNotFound(err) {
				return err
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("Starting context '%s'", context.Name))

			cmdArgs := []string{}
			if context.JVMOptions != nil && len(context.JVMOptions) > 0 {
				cmdArgs = append(cmdArgs, context.JVMOptions...)
			}

			if additionalJvmOptions != nil && len(additionalJvmOptions) > 0 {
				cmdArgs = append(cmdArgs, additionalJvmOptions...)
			}
			cmdArgs = append(cmdArgs, "-jar", context.War)

			if context.Options != nil {
				cmdArgs = append(cmdArgs, context.Options...)
			}

			c := exec.Command("java", cmdArgs...)

			if context.JavaHome != "" {
				c.Env = append(c.Env, fmt.Sprintf("JAVA_HOME=%s", context.JavaHome))
			}
			if context.JenkinsHome != "" {
				c.Env = append(c.Env, fmt.Sprintf("JENKINS_HOME=%s", context.JenkinsHome))
			}
			c.Env = append(c.Env, "JENKINS_HA=false")

			err := pkg.DisplayCommandOutput(c)
			context.Pid = c.Process.Pid
			config.Config.Jenkins.UpdateContext(context)
			config.Save()
			return err
		},
	}
	startCommand.Flags().StringArrayVar(&additionalJvmOptions, "jvm-option", nil, "The JVM option to pass when starting Jenkins. Use this option multiple times for many options")
	command.AddCommand(startCommand)

	stopCommand := &cobra.Command{
		Use:   "stop",
		Short: "Stops a Jenkins context",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToStop := ""
			if len(args) == 1 {
				contextToStop = args[0]
			} else if config.Config.Jenkins.CurrentContext != "" {
				contextToStop = config.Config.Jenkins.CurrentContext
			} else {
				return errors.New("No Jenkins context specified, neither a current one is set")
			}

			if c, err := config.Config.Jenkins.FindContextByName(contextToStop); config.IsJenkinsContextNotFound(err) {
				return err
			} else if c.Pid == -1 {
				return errors.New(fmt.Sprintf("The context %s is not started", contextToStop))
			} else {
				context = c
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(fmt.Sprintf("Stopping context '%s'", context.Name))
			err := syscall.Kill(context.Pid, syscall.SIGTERM)

			if err == nil {
				context.Pid = -1
				config.Config.Jenkins.UpdateContext(context)
				config.Save()
			}

			return err
		},
	}
	stopCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	command.AddCommand(stopCommand)

	passwordCommand := &cobra.Command{
		Use:   "admin-password",
		Short: "Gets the initial admin password of Jenkins",
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			contextToGetPasswordFor := ""
			if len(args) == 1 {
				contextToGetPasswordFor = args[0]
			} else if config.Config.Jenkins.CurrentContext != "" {
				contextToGetPasswordFor = config.Config.Jenkins.CurrentContext
			} else {
				return errors.New("No Jenkins context specified, neither a current one is set")
			}

			if c, err := config.Config.Jenkins.FindContextByName(contextToGetPasswordFor); config.IsJenkinsContextNotFound(err) {
				return err
			} else {
				context = c
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			bytes, err := os.ReadFile(filepath.Join(context.JenkinsHome, "secrets", "initialAdminPassword"))
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
			return nil
		},
	}
	command.AddCommand(passwordCommand)

	return command
}
