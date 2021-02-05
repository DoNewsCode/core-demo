package cmd

import (
	"github.com/DoNewsCode/skeleton/cmd/internal"
)

// Execute executes the root command.
func Execute() error {
	rootCmd, c := internal.Bootstrap()
	defer c.Shutdown()

	// setup command graph
	rootCmd.AddCommand(NewVersionCommand(c))
	rootCmd.AddCommand(NewServeCommand(c))
	rootCmd.AddCommand(NewMigrateCommand(c))
	rootCmd.AddCommand(NewSeedCommand(c))
	rootCmd.AddCommand(NewSeedRedisCommand(c))

	// run
	return rootCmd.Execute()
}
