package internal

import (
	"math/rand"
	"os"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/spf13/cobra"
)

func Bootstrap() (*cobra.Command, *core.C) {

	rand.Seed(time.Now().UnixNano())

	var cfgPath string

	rootCmd := &cobra.Command{
		Use:   "Skeleton",
		Short: "A Pragmatic and Opinionated Go Application",
		Long:  `Skeleton provides a starting point to write 12-factor Go Applications.`,
	}

	// Determine config path from commandline
	rootCmd.PersistentFlags().StringVar(
		&cfgPath,
		"config",
		"./config/config.yaml",
		"config file (default is ./config/config.yaml)",
	)
	_ = rootCmd.PersistentFlags().Parse(os.Args[1:])

	// setup core with config file path
	c := core.Default(core.WithYamlFile(cfgPath))

	// setup global dependencies
	provide(c)

	// register global modules
	register(c)

	// add commands from modules
	c.ApplyRootCommand(rootCmd)

	return rootCmd, c
}
