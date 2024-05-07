package cmd

import (
	"fmt"
	"giftcard/internal/adaptor/postgres"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
	Long:  `Commands to manage database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Applying models migration done.")
		if err := postgres.MigrateModels(); err != nil {
			fmt.Println("error: ", err)
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
