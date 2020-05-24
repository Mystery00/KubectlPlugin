package cmd

import (
	"KubectlPlugin/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "安装工具到PATH中。",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := exec.LookPath(os.Args[0])
		if err != nil {
			panic(err)
		}
		path, err := filepath.Abs(file)
		if err != nil {
			panic(err)
		}
		currentFileName := filepath.Base(path)
		// 创建文件
		_, err = os.Create(installPath)
		if err != nil {
			panic(err)
		}
		err = os.Rename(currentFileName, installPath)
		if err != nil {
			panic(err)
		}
		fmt.Println(utils.INFO + " 工具安装成功！")
		fmt.Println(utils.INFO + " 可以通过 " + utils.Green("k8s") + " 直接运行脚本")
	},
}
