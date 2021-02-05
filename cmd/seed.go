package cmd

import (
	"fmt"
	"os"

	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/DoNewsCode/std/pkg/core"
	"github.com/DoNewsCode/std/pkg/otgorm"
	"github.com/spf13/cobra"
)

func NewSeedCommand(c *core.C) *cobra.Command {
	var force bool
	var seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "seed the database",
		Long:  `use the provided seeds to bootstrap fake data in database`,
		Run: func(cmd *cobra.Command, args []string) {
			c.Invoke(func(env contract.Env) {
				if env.IsProduction() && !force {
					c.Err(fmt.Errorf("seeding in production requires force flag to be set"))
					os.Exit(1)
				}
			})

			seeds := collectSeeds(c)

			if err := seeds.Seed(); err != nil {
				c.Err(fmt.Errorf("seed failed: %w", err))
				os.Exit(1)
			}

			c.Info("seeding successfully completed")
		},
	}
	seedCmd.Flags().BoolVarP(&force, "force", "f", false, "seeding in production requires force flag to be set")

	return seedCmd
}

func collectSeeds(c *core.C) otgorm.Seeds {
	var seeds otgorm.Seeds
	for _, f := range c.SeedProviders {
		seeds.Collection = append(seeds.Collection, f()...)
	}
	seeds.Logger = c
	err := c.Populate(&seeds.Db)
	if err != nil {
		c.Err(err)
		os.Exit(1)
	}
	return seeds
}
