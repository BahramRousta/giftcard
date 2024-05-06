package cmd

import (
	"fmt"
	"giftcard/app"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "giftcard",
	Short: "Base command for gift card",
	Long:  "Base command for gift card",
	Run: func(cmd *cobra.Command, args []string) {
		app.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
