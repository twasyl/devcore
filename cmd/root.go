package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devcore",
	Short: "devcore provides developers/development utilities",
}

func init() {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install devcore in your environment",
		Run: func(cmd *cobra.Command, args []string) {
			_, file, _, _ := runtime.Caller(0)
			fmt.Println(file)
		},
	}

	rootCmd.AddCommand(cmd)
}

// Execute will execute the `devcore command`
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
