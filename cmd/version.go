package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"KubectlPlugin/utils"
)

const version = "V1.1.4"

const versionTpl = utils.INFO + " k8s服务连接脚本 {{\"[\"|red}}{{.|red}}{{\"]\"|red}}  -- Made with {{\"♥\"|red}} by {{\"Mystery0\"|blue}}"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本号。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.Parse(versionTpl, version))
	},
}
