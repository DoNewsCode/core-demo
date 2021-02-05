package cmd

import (
	"fmt"
	"os"

	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/DoNewsCode/std/pkg/core"
	"github.com/DoNewsCode/std/pkg/otgorm"
	"github.com/spf13/cobra"
)

func NewMigrateCommand(c *core.C) *cobra.Command {
	var (
		force      bool
		rollbackId string
	)
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate gorm tables",
		Long:  `Run all gorm table migrations.`,
		Run: func(cmd *cobra.Command, args []string) {
			c.Invoke(func(env contract.Env) {
				if env.IsProduction() && !force {
					c.Err(fmt.Errorf("migrations and rollback in production requires force flag to be set"))
					os.Exit(1)
				}
			})

			migrations := collectMigrations(c)

			if rollbackId != "" {
				if err := migrations.Rollback(rollbackId); err != nil {
					c.Err(fmt.Errorf("unable to rollback: %w", err))
					os.Exit(1)
				}

				c.Info("rollback successfully completed")
				return
			}

			if err := migrations.Migrate(); err != nil {
				c.Err(fmt.Errorf("unable to migrate: %w", err))
				os.Exit(1)
			}

			c.Info("migration successfully completed")
		},
	}
	migrateCmd.Flags().BoolVarP(&force, "force", "f", false, "migrations and rollback in production requires force flag to be set")
	migrateCmd.Flags().StringVarP(&rollbackId, "rollback", "r", "", "rollback to the given migration id")
	migrateCmd.Flag("rollback").NoOptDefVal = "-1"
	return migrateCmd
}

func collectMigrations(c *core.C) otgorm.Migrations {
	var migrations otgorm.Migrations
	for _, p := range c.MigrationProviders {
		migrations.Collection = append(migrations.Collection, p()...)
	}
	err := c.Populate(&migrations.Db)
	c.CheckErr(err)

	return migrations
}
