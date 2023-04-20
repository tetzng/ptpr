/*
Copyright © 2023 Teppei Taguchi tetzng.tt@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ptpr",
	Short: "create PR template from Pivotal Tracker",
	Long:  `create Pull Request template from your project's Pivotal Tracker.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
