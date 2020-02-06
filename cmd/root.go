package cmd

import (
	"fmt"
	"os"

	"github.com/chrisdail/pile/gitver"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var rootDir string

var rootCmd = &cobra.Command{
	Use:     "pile",
	Version: version,
	Short:   "Simple docker container builder",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		gitver.SetWorkingDir(rootDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "", "Root workspace directory (defaults to git root)")
}
