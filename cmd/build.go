package cmd

import (
	"fmt"
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

		var lastBuildErr error
		var images []string
		for _, project := range projects {
			buildImage, err := piler.Build(&project)
			if err != nil {
				log.Println(err)
				lastBuildErr = err
			} else if buildImage.FullyQualifiedImage != "" {
				images = append(images, buildImage.FullyQualifiedImage)
			}
		}

		if lastBuildErr != nil {
			return lastBuildErr
		}

		// If only a single image was built, print out only the image name
		// This allows using the stdout as an input to other commands
		if len(images) == 1 {
			fmt.Print(images[0])
		} else if len(images) > 1 {
			for _, image := range images {
				fmt.Println(image)
			}
		}
		return nil
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
