package cmd

import (
	"KubectlPlugin/utils"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新工具版本。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.Parse(versionTpl, version))
	},
}
