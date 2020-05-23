package cmd

import (
	"KubectlPlugin/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装工具到PATH中。",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Rename(k8sFileName, installPath)
		if err != nil {
			panic(err)
		}
		err = os.RemoveAll(fileName)
		if err != nil {
			panic(err)
		}
		fmt.Println(utils.INFO + " 工具安装成功！")
		fmt.Println(utils.INFO + " 可以通过 " + utils.Green("k8s") + " 直接运行脚本")
	},
}
