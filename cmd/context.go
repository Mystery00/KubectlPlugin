package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"KubectlPlugin/mritd"
	"KubectlPlugin/utils"
)

var currentContext string

const contextTpl = utils.INFO_TPL + " 当前集群环境为 {{.|red}}"
const changeContextTpl = utils.INFO_TPL + " 集群环境切换为： {{.|cyan}}"

func init() {
	rootCmd.AddCommand(contextCmd)
	rootCmd.PersistentFlags().StringVarP(&currentContext, "context", "c", "", "指定本次运行的集群环境")
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "仅切换集群环境。",
	Long:  `切换集群环境，并且应用到 kubectl 的配置文件中。`,
	Run: func(cmd *cobra.Command, args []string) {
		//快速修改上下文
		if currentContext == "" {
			selectContext()
		} else {
			setContext(currentContext)
		}
	},
}

func selectContext() {
	fmt.Println()
	cmdStr := fmt.Sprintln(`kubectl config get-contexts | grep -v 'NAME' | sed 's/^.//g' | awk '{print $1}'`)
	result := utils.Cmd("sh", "-c", cmdStr)
	contextList := strings.Fields(result)
	maxLength := 28

	if contextList == nil {
		fmt.Println("没有集群环境配置，退出中...")
		os.Exit(0)
	}
	length := maxLength + 16
	cfg := &mritd.SelectConfig{
		ActiveTpl:       utils.LINE + " {{\"↣  Switch to \"|blue}}{{.|lLength " + mritd.ToString(maxLength) + "|magenta}}{{\"  ↢\"|blue}} " + utils.LINE,
		InactiveTpl:     utils.LINE + " {{\"   Switch to \"|white}}{{.|lLength " + mritd.ToString(maxLength) + "|cyan}}    " + utils.LINE,
		SelectedTpl:     utils.INFO_TPL + " 集群环境切换为： {{.|cyan}}",
		DisPlaySize:     len(contextList),
		SelectHeaderTpl: utils.LINE + " {{\"Select Context:\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		SelectPromptTpl: utils.LINE + " {{\"Use the arrow keys to navigate: ↓ ↑ → ←\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		ShowBorder:      true,
		ShowWidth:       length,
	}
	s := &mritd.Select{
		Items:  contextList,
		Config: cfg,
	}
	context := contextList[s.Run()]
	setContextOnly(context)
	fmt.Println()
}

func setContextOnly(context string) {
	utils.Cmd("kubectl", "config", "use-context", context)
	currentContext = context
}

func setContext(context string) {
	setContextOnly(context)
	fmt.Println(utils.Parse(changeContextTpl, context))
}

func getCurrentContext() string {
	context := strings.TrimSpace(utils.Cmd("kubectl", "config", "current-context"))
	fmt.Println(utils.Parse(contextTpl, context))
	currentContext = context
	return context
}
