package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"KubectlPlugin/mritd"
	"KubectlPlugin/utils"
)

var currentNamespace string

const namespaceTpl = utils.INFO_TPL + " 检测到默认命名空间 [{{.|cyan}}]"
const changeNamespaceTpl = utils.INFO_TPL + " 集群环境 [{{.context|blue}}] 默认命名空间切换为： {{.namespace|cyan}}"

func init() {
	rootCmd.AddCommand(nameSpaceCmd)
	rootCmd.PersistentFlags().StringVarP(&currentNamespace, "namespace", "n", "", "指定本次运行的命名空间")
}

var nameSpaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "仅切换命名空间。",
	Long:  `将命名空间设置为当前集群环境的默认命名空间。`,
	Run: func(cmd *cobra.Command, args []string) {
		//快速修改命名空间
		if currentContext == "" {
			//没有设置上下文，查询
			getCurrentContext()
		}
		//快速修改命名空间
		if currentNamespace == "" {
			selectNamespace()
		} else {
			setNamespace(currentContext, currentNamespace)
		}
	},
}

func detectNamespace() {
	b, namespace := checkNamespace()
	if b {
		fmt.Println(utils.Parse(namespaceTpl, namespace))
	} else {
		namespace = selectNamespace()
	}
	currentNamespace = namespace
}
func checkNamespace() (bool, string) {
	cmdStr := fmt.Sprintf(`kubectl config get-contexts %s | grep -v 'NAME' | sed 's/^.//g' | awk '{print $1,$4}'`, currentContext)
	result := utils.Cmd("sh", "-c", cmdStr)
	res := strings.Fields(result)
	if len(res) == 2 {
		return true, res[1]
	} else {
		return false, ""
	}
}

func selectNamespace() string {
	fmt.Println()
	result := utils.Cmd("kubectl", "get", "namespace", "-o", "json")
	//获取命名空间列表
	namespaces := gjson.Get(result, "items").Array()
	//存储解析后的命名空间
	var namespaceList []string
	maxLength := 28

	for _, namespace := range namespaces {
		ns := gjson.Get(namespace.Raw, "metadata.name").String()
		namespaceList = append(namespaceList, ns)
		maxLength = utils.Max(maxLength, len(ns))
	}

	if namespaceList == nil {
		fmt.Println("没有命名空间，退出中...")
		os.Exit(0)
	}
	length := maxLength + 16
	cfg := &mritd.SelectConfig{
		ActiveTpl:       utils.LINE + " {{\"↣  Switch to \"|blue}}{{.|lLength " + mritd.ToString(maxLength) + "|magenta}}{{\"  ↢\"|blue}} " + utils.LINE,
		InactiveTpl:     utils.LINE + " {{\"   Switch to \"|white}}{{.|lLength " + mritd.ToString(maxLength) + "|cyan}}    " + utils.LINE,
		SelectedTpl:     utils.INFO_TPL + " 该集群环境 [{{\"" + currentContext + "\"|blue}}] 默认命名空间切换为： {{.|cyan}}",
		DisPlaySize:     len(namespaceList),
		SelectHeaderTpl: utils.LINE + " {{\"Select Namespace:\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		SelectPromptTpl: utils.LINE + " {{\"Use the arrow keys to navigate: ↓ ↑ → ←\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		ShowBorder:      true,
		ShowWidth:       length,
	}
	s := &mritd.Select{
		Items:  namespaceList,
		Config: cfg,
	}
	namespace := namespaceList[s.Run()]
	setNamespaceOnly(currentContext, namespace)
	fmt.Println()
	return namespace
}

func setNamespaceOnly(context string, namespace string) {
	utils.Cmd("kubectl", "config", "set-context", context, "--namespace="+namespace)
	currentNamespace = namespace
}

func setNamespace(context string, namespace string) {
	setNamespaceOnly(context, namespace)
	fmt.Println(utils.Parse(changeNamespaceTpl, struct {
		context   string
		namespace string
	}{context, namespace}))
}
