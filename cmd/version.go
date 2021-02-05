package cmd

import (
	"fmt"

	"github.com/DoNewsCode/std/pkg/core"
	"github.com/spf13/cobra"
)

func NewVersionCommand(c *core.C) *cobra.Command {
	var greeting string
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Log the version number of this app",
		Long:  `Under semantic versioning, version numbers and the way they change convey meaning about the underlying code and what has been modified from one version to the next.`,
		Run: func(cmd *cobra.Command, args []string) {
			version := c.String("version")
			c.Info(fmt.Sprintf("%s, the version is %s", greeting, version))
		},
	}
	versionCmd.Flags().StringVar(&greeting, "greeting", "hello", "a greeting message")
	return versionCmd
}
