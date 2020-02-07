package cmd

import (
	"fmt"
	"log"

	"github.com/chrisdail/pile/core"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info [projects...]",
	Short: "Lists info about projects",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Workspace: %s", core.Workspace.Dir)
		projects, err := core.Workspace.ProjectsFromArgs(args)
		if err != nil {
			return err
		}

		log.Println("Projects")
		for _, project := range projects {
			fmt.Printf("- %s\n", project.Dir)
			if project.CanBuild {
				fmt.Printf("  Builds: %s\n  Pushes As: %s\n",
					project.Image,
					project.FullyQualifiedImage,
				)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
