package cmd

import (
	"KubectlPlugin/utils"
	"archive/zip"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const user = "Mystery00"
const project = "KubectlPlugin"
const fileName = "k8s.zip"
const k8sFileName = "k8s"
const installPath = "/usr/local/bin/k8s"
const apiUrl = "https://api.github.com/repos/" + user + "/" + project + "/releases/latest"
const downloadUrl = "https://github.com/" + user + "/" + project + "/releases/download/%s/" + fileName

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新工具版本。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utils.INFO + " 正在检查更新...")
		cmdStr := `wget -qO- -t1 -T2 "` + apiUrl + `" | grep "tag_name" | head -n 1 | awk -F ":" '{print $2}' | sed 's/\"//g;s/,//g;s/ //g'`
		latestVersion := strings.TrimSpace(utils.Cmd("sh", "-c", cmdStr))
		if latestVersion != version {
			fmt.Println(utils.INFO + " 检测到新版本 [" + latestVersion + "] ，正在下载...")
			downloadFile(latestVersion)
			copyFile()
			fmt.Println("工具已更新为最新版本 [" + utils.Red(latestVersion) + "] !(注意：因为更新方式为直接覆盖当前运行的脚本，所以可能下面会提示一些报错，无视即可)")
		} else {
			fmt.Println(utils.WARN + " 当前已经是最新版本！")
		}
	},
}

func downloadFile(latestVersion string) {
	//版本不一致，更新版本
	resp, err := http.Get(fmt.Sprintf(downloadUrl, latestVersion))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}

func copyFile() {
	err := unzip(fileName, ".")
	if err != nil {
		panic(err)
	}
	err = os.Rename(k8sFileName, installPath)
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
