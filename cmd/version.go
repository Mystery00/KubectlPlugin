package cmd

import (
	"KubectlPlugin/utils"
	"fmt"
	"github.com/spf13/cobra"
)

const version = "1.1.0"

const versionTpl = utils.INFO + " k8s服务连接脚本 {{\"[v\"|red}}{{.|red}}{{\"]\"|red}}  -- Made with ♥ by Mystery0"

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
