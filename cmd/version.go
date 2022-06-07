package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	c "io.twasyl/devcore/pkg/config"
)

func init() {
	rootCmd.AddCommand(buildVersionCmd())
}

func buildVersionCmd() *cobra.Command {
	verbose := false
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display devcore version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("devcore version 1")

			if verbose {
				fmt.Println("- Supported tools for installation with their default version:")
				for _, tool := range c.DefaultSupportedTools {
					fmt.Printf("  - %s %s\n", tool.Name, tool.DefaultVersion)
				}
				fmt.Println(`- Docker compose contexts can be:
  - added
  - deleted
  - listed
  - opened in the file explorer
- The Jenkins CLI can be:
  - downloaded
  - executed
- git projects can be:
  - cloned
  - checked if they have already been cloned
- Helm charts can be:
  - installed
  - uninstalled
- With kind, clusters can be:
  - created with a Kubernetes dashboard (as well as getting the connection token)
  - deleted
- The helper allows to properly uninstall docker`)
			}
		},
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Displays more info concerning the devcore version")
	return cmd
}
