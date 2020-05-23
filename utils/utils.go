package utils

import (
	"KubectlPlugin/mritd"
	"github.com/mritd/promptx/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const SpaceSeparatorLength = 1

const (
	CENTER        = "═"
	LINE          = "║"
	LEFT_TOP      = "╔"
	CENTER_TOP    = "╦"
	RIGHT_TOP     = "╗"
	LEFT_CENTER   = "╠"
	CENTER_CENTER = "╬"
	RIGHT_CENTER  = "╣"
	LEFT_BOTTOM   = "╚"
	CENTER_BOTTOM = "╩"
	RIGHT_BOTTOM  = "╝"
)

const (
	INFO_TPL  = "{{\"[信息]\"|green}}"
	ERROR_TPL = "{{\"[错误]\"|red}}"
	WARN_TPL  = "{{\"[注意]\"|yellow}}"
)
const (
	INFO  = "\033[32m[信息]\033[0m"
	ERROR = "\033[31m[错误]\033[0m"
	WARN  = "\033[33m[注意]\033[0m"
)

func Cmd(name string, arg ...string) string {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	b, err := cmd.Output()
	if err != nil {
		os.Exit(1)
	}
	return string(b)
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func PrintCenter(str string, length int) string {
	stringLength := len(str)
	leftLength := (length - stringLength) / 2
	rightLength := length - stringLength - leftLength
	return strings.Repeat(" ", leftLength) + str + strings.Repeat(" ", rightLength)
}

func Parse(tpl string, data interface{}) string {
	res, err := template.New("").Funcs(mritd.FuncMap).Parse(tpl)
	if err != nil {
		os.Exit(1)
	}
	return string(utils.Render(res, data))
}

func DownloadFile(url string, fileName string) {
	resp, err := http.Get(url)
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
