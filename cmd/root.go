package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/chrisdail/pile/core"

	"github.com/chrisdail/pile/gitver"

	"github.com/spf13/cobra"
)

var rootDir string

var rootCmd = &cobra.Command{
	Use:          "pile",
	Short:        "Simple docker container builder",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !filepath.IsAbs(rootDir) {
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}

			rootDir = filepath.Join(workingDir, rootDir)
		}
		gitver.SetWorkingDir(rootDir)
		return core.Workspace.SetDir(rootDir)
	},
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, core.ErrorTestsFailed) {
			os.Exit(127)
		} else {
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "", "Root workspace directory (defaults to git root)")
}
