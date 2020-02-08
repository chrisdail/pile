package cmd

import (
	"log"

	"github.com/chrisdail/pile/core"
	"github.com/spf13/cobra"
)

var piler = &core.Piler{}

var buildCmd = &cobra.Command{
	Use:   "build [projects...]",
	Short: "Builds a set of projects",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := core.Workspace.ProjectsFromArgs(args)
		if err != nil {
			return err
		}

		for _, project := range projects {
			_, err = piler.Build(&project)
			if err != nil {
				log.Println(err)
			}
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&piler.Force, "force", "f", false,
		"Forces a rebuild of the container even if one exists already built")
	buildCmd.Flags().BoolVar(&piler.SkipPush, "skip-push", false,
		"Skips pushing the container to the remote registry")
	buildCmd.Flags().BoolVar(&piler.SkipTests, "skip-tests", false,
		"Skips running the tests")
}
