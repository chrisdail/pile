package cmd

import (
	"fmt"
	"os"

	"github.com/chrisdail/pile/core"

	"github.com/chrisdail/pile/gitver"

	"github.com/spf13/cobra"
)

var rootDir string

var rootCmd = &cobra.Command{
	Use:   "pile",
	Short: "Simple docker container builder",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		gitver.SetWorkingDir(rootDir)
		return core.Workspace.SetDir(rootDir)
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "", "Root workspace directory (defaults to git root)")
}
