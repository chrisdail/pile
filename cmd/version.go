package cmd

import (
	"fmt"

	"github.com/chrisdail/pile/gitver"
	"github.com/spf13/cobra"
)

var versionTemplate string

var versionCmd = &cobra.Command{
	Use:   "version [projects...]",
	Short: "Generates a git-based version for projects relative to the git root",
	Example: `  pile version
  pile version app/backend ui
  pile version -t "{{.Branch}}.{{.Hash}}"`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := gitver.ForProjects(args)
		if err != nil {
			return err
		}

		formatted, err := version.FormatTemplate(versionTemplate)
		if err != nil {
			return err
		}
		fmt.Println(formatted)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().StringVarP(&versionTemplate, "template", "t", gitver.DefaultTemplate,
		"Text template used to format the version")
}
