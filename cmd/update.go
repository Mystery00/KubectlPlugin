package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"KubectlPlugin/utils"
)

const user = "Mystery00"
const project = "KubectlPlugin"
const installPath = "/usr/local/bin/k8s"
const tempDir = "/tmp/k8s/"
const apiUrl = "https://api.github.com/repos/" + user + "/" + project + "/releases/latest"
const downloadUrl = "https://github.com/" + user + "/" + project + "/releases/download/%s/"

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新工具版本。",
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			panic(err)
		}
		var fileName string
		var k8sFileName string
		switch runtime.GOOS {
		case "darwin":
			fileName = "k8s-mac.zip"
			k8sFileName = "k8s-mac"
		case "windows":
			fmt.Println(utils.ERROR + " 暂不支持 Windows")
			os.Exit(0)
		case "linux":
			fileName = "k8s-linux.zip"
			k8sFileName = "k8s-linux"
		default:
			fmt.Println(utils.ERROR + " 未知操作系统")
			os.Exit(0)
		}
		fmt.Println(utils.INFO + " 正在检查更新...")
		cmdStr := `wget -qO- -t1 -T2 "` + apiUrl + `" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
		latestVersion := strings.TrimSpace(utils.Cmd("sh", "-c", cmdStr))
		if latestVersion != version {
			//版本不一致，更新版本
			fmt.Println(utils.INFO + " 检测到新版本 [" + latestVersion + "] ，正在下载...")
			utils.DownloadFile(fmt.Sprintf(downloadUrl+fileName, latestVersion), tempDir+fileName)
			copyFile(tempDir+fileName, k8sFileName)
			fmt.Println("工具已更新为最新版本 [" + utils.Red(latestVersion) + "] !(注意：因为更新方式为直接覆盖当前运行的脚本，所以可能下面会提示一些报错，无视即可)")
		} else {
			fmt.Println(utils.WARN + " 当前已经是最新版本！")
		}
	},
}

func copyFile(fileName string, k8sFileName string) {
	err := unzip(fileName, tempDir)
	if err != nil {
		panic(err)
	}
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	currentFileName := filepath.Base(path)
	err = os.Rename(tempDir+k8sFileName, currentFileName)
	if err != nil {
		panic(err)
	}
	err = os.RemoveAll(fileName)
	if err != nil {
		panic(err)
	}
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
