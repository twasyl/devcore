package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"io.twasyl/devcore/pkg/config"
	pkg "io.twasyl/devcore/pkg/utils"
)

func init() {
	var projectCmd = &cobra.Command{
		Use:   "project",
		Short: "Utilities to manage GitHub projects",
	}
	projectCmd.AddCommand(buildProjectExistCmd())
	projectCmd.AddCommand(buildProjectCloneCmd())
	projectCmd.AddCommand(buildProjectUpdateCmd())
	rootCmd.AddCommand(projectCmd)

}

func buildProjectExistCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exist",
		Short: "Check if a project exist locally",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			project := filepath.Join(config.Config.ProjectsDir, args[0])
			if _, err := os.Stat(project); err == nil {
				fmt.Printf("Project %s exists in %s\n", args[0], config.Config.ProjectsDir)
			} else if os.IsNotExist(err) {
				fmt.Printf("Project %s does not exist in %s\n", args[0], config.Config.ProjectsDir)
			} else {
				return err
			}
			return nil
		},
	}
	return cmd
}

func buildProjectCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a project from GitHub",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := pkg.ExecCommandInDir(config.Config.ProjectsDir, "git", "clone", fmt.Sprintf("git@github.com:%s.git", args[0]))
			fmt.Println(out)
			return err
		},
	}
	return cmd
}

func buildProjectUpdateCmd() *cobra.Command {
	var forceOnError = false
	var verbose = false

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a git project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Updating project %s\n", args[0])
			gitProjectDir := path.Join(config.Config.ProjectsDir, args[0])
			out, err := pkg.ExecCommandInDir(gitProjectDir, "git", "pull")

			if verbose {
				fmt.Println(out)
			}
			if err != nil && !forceOnError {
				return err
			}

			// Listing git subprojects and update them
			dir, err := os.Open((gitProjectDir))
			if err != nil {
				return err
			}

			subProjects, err := dir.ReadDir(-1)
			if err != nil {
				return err
			}

			for _, subProject := range subProjects {
				if subProject.IsDir() {
					if _, err := os.Stat(path.Join(gitProjectDir, subProject.Name(), ".git")); !os.IsNotExist(err) {
						fmt.Printf("Updating subproject %s\n", subProject.Name())
						out, err := pkg.ExecCommandInDir(path.Join(gitProjectDir, subProject.Name()), "git", "pull")

						if verbose {
							fmt.Println(out)
						}
						if err != nil && !forceOnError {
							return err
						}
					}
				}
			}
			return err
		},
	}
	cmd.Flags().BoolVarP(&forceOnError, "force", "f", false, "Force updating subprojects when an error occurs")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output of update")
	return cmd
}
