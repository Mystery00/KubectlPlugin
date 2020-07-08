package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"

	"KubectlPlugin/mritd"
	"KubectlPlugin/utils"
)

var rootCmd = &cobra.Command{
	Use:   "k8s",
	Short: "k8s服务连接工具",
	Long:  `一个简单的k8s服务连接工具，使用Go语言编写。仓库地址： https://github.com/Mystery00/KubectlPlugin`,
	Run: func(cmd *cobra.Command, args []string) {
		doAction()
	},
}

type options int

const (
	OptionEnterPod options = iota
	OptionQueryPort
	OptionSwitchContext
	OptionSwitchNamespace
)

type Pod struct {
	app        string
	branch     string
	name       string
	containers []Container
}

type Container struct {
	name string
}

func doAction() {
	if currentContext == "" {
		fmt.Println("正在查询当前集群环境...")
		getCurrentContext()
	}
	if currentNamespace == "" {
		detectNamespace()
	}
	for true {
		b, pod := handleOption(selectPod())
		if b {
			input := shouldContinue()
			if strings.ToLower(input) == "y" {
				continue
			}
			os.Exit(0)
		}
		fmt.Println()
		//检查容器名称是否正确
		var isContainerNameInvalid = true
		var containerName string
		i := strings.Index(pod.name, "-xylinkapp")
		if i != -1 {
			containerName = pod.name[0:i]
			for _, container := range pod.containers {
				if containerName == container.name {
					isContainerNameInvalid = false
					break
				}
			}
		}
		if isContainerNameInvalid {
			fmt.Println(utils.WARN + " 容器名称自动检测失败，请手动选择...")
			container := selectContainer(pod.containers)
			containerName = container.name
		}
		cmd := exec.Command("kubectl", "exec", "-it", "-n", currentNamespace, pod.name, "-c", containerName, "--", "/bin/bash")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		connectTpl := utils.INFO_TPL + " 正在连接到 [{{.|cyan}}] 的Shell..."
		fmt.Println(utils.Parse(connectTpl, pod.name))
		fmt.Println(strings.Repeat(utils.CENTER, 40) + " 输出开始 " + strings.Repeat(utils.CENTER, 40))
		_ = cmd.Run()
		fmt.Println(strings.Repeat(utils.CENTER, 40) + " 输出结束 " + strings.Repeat(utils.CENTER, 40))
		fmt.Println(utils.INFO + " Shell已退出...")
		fmt.Println()
		if fastIndex != -1 {
			//清除快速序号
			fastIndex = -1
		}
		input := shouldContinue()
		if strings.ToLower(input) == "y" {
			continue
		}
		os.Exit(0)
	}
}

func shouldContinue() string {
	var input string
	fmt.Print(utils.INFO + " 是否需要继续操作？[y/N] ")
	_, _ = fmt.Scanln(&input)
	return input
}

func handleOption(option options, pod Pod) (bool, Pod) {
	switch option {
	case OptionEnterPod:
		return false, pod
	case OptionQueryPort:
		queryPort()
	case OptionSwitchContext:
		selectContext()
	case OptionSwitchNamespace:
		selectNamespace()
	default:
		fmt.Println(utils.ERROR + " 未知的指令")
	}
	return true, Pod{}
}

func selectPod() (options, Pod) {
	fmt.Println()
	fmt.Println(utils.INFO + " 正在查询Pod信息...")
	result := utils.Cmd("kubectl", "get", "pod", "--namespace="+currentNamespace, "-o", "json")
	//获取pod列表-json
	pods := gjson.Get(result, "items").Array()
	//存储解析后的pod列表
	var originPodList []Pod
	var podList []Pod

	//第一列最大宽度
	var maxLength1 = 4
	//第二列最大宽度
	var maxLength2 = 3
	//第三列最大宽度
	var maxLength3 = 3
	for _, podRaw := range pods {
		app := gjson.Get(podRaw.Raw, "metadata.labels.app")
		branch := gjson.Get(podRaw.Raw, "metadata.annotations.branch")
		name := gjson.Get(podRaw.Raw, "metadata.name")
		var containerList []Container
		for _, containerRaw := range gjson.Get(podRaw.Raw, "spec.containers").Array() {
			containerList = append(containerList, Container{name: gjson.Get(containerRaw.Raw, "name").String()})
		}
		pod := Pod{app.String(), branch.String(), name.String(), containerList}
		if pod.app == "" {
			pod.app = "[not app]"
		}
		if pod.branch == "" {
			pod.branch = "[no branch]"
		}
		originPodList = append(originPodList, pod)
		maxLength1 = utils.Max(maxLength1, len(pod.app))
		maxLength2 = utils.Max(maxLength2, len(pod.branch))
		maxLength3 = utils.Max(maxLength3, len(pod.name))
	}
	if originPodList == nil {
		fmt.Println(utils.ERROR + " ☟☟️ 该命名空间中没有Pod，请切换命名空间...☟☟")
		return OptionSwitchNamespace, Pod{}
	}
	if isFast || len(originPodList) > 20 {
		if !isFast {
			//选项数量多于20个，询问是否通过关键词进行过滤
			fmt.Print("Pod数量过多，请输入过滤关键词：(默认-不过滤)")
		}
		var input string
		for true {
			if !isFast {
				_, _ = fmt.Scanln(&input)
				if input == "" {
					//没有输入文本，跳过过滤
					//拷贝Pod列表
					podList = originPodList
					//跳出循环
					break
				}
			} else {
				input = fastName
			}
			//重置最大宽度
			maxLength1 = 4
			maxLength2 = 3
			maxLength3 = 3
			//清空旧的列表
			podList = nil
			for _, pod := range originPodList {
				if strings.Contains(pod.name, input) {
					podList = append(podList, pod)
					maxLength1 = utils.Max(maxLength1, len(pod.app))
					maxLength2 = utils.Max(maxLength2, len(pod.branch))
					maxLength3 = utils.Max(maxLength3, len(pod.name))
				}
			}
			if podList != nil {
				break
			} else {
				fmt.Print(utils.WARN + " 没有Pod满足过滤条件，请重新输入过滤关键词：")
			}
		}
	} else {
		podList = originPodList
	}
	maxOption := len(podList)

	//|1+4+1|1+max1+1|1+max2+1|1+max3+1|
	length := (0) + utils.SpaceSeparatorLength + 4 + utils.SpaceSeparatorLength + (1) + utils.SpaceSeparatorLength + (11 + maxLength1) + utils.SpaceSeparatorLength + (1) + utils.SpaceSeparatorLength + (8 + maxLength2) + utils.SpaceSeparatorLength + (1) + utils.SpaceSeparatorLength + (maxLength3) + utils.SpaceSeparatorLength + (0)
	specificOptionLength := utils.SpaceSeparatorLength + (11 + maxLength1) + utils.SpaceSeparatorLength + (1) + utils.SpaceSeparatorLength + (8 + maxLength2) + utils.SpaceSeparatorLength + (1) + utils.SpaceSeparatorLength + (maxLength3) + utils.SpaceSeparatorLength

	specificOptionTpl := utils.LINE + ` {{.INDEX|lLength 4|cyan}} ` + utils.LINE + ` {{.OPTION|lLength ` + mritd.ToString(specificOptionLength-2) + `}} ` + utils.LINE
	specificOptionWithInfoTpl := utils.LINE + ` {{.INDEX|lLength 4|cyan}} ` + utils.LINE + ` {{.OPTION|lLength ` + mritd.ToString(maxLength1+maxLength2) + `}} Current: {{.INFO|lLength ` + mritd.ToString(specificOptionLength-12-maxLength1-maxLength2) + `|red}} ` + utils.LINE
	optionTpl := utils.LINE + ` {{.INDEX|lLength 4|cyan}} ` + utils.LINE + ` Connect to {{.OPTION|lLength ` + mritd.ToString(maxLength1) + `|magenta}} ` + utils.LINE + ` Branch: {{.BRANCH|lLength ` + mritd.ToString(maxLength2) + `|green}} ` + utils.LINE + ` {{.NAME|lLength ` + mritd.ToString(maxLength3) + `|blue}} ` + utils.LINE

	fmt.Println(utils.LEFT_TOP + strings.Repeat(utils.CENTER, length) + utils.RIGHT_TOP)
	fmt.Println(utils.LINE + utils.PrintCenter("Please select the operation option to be performed", length) + utils.LINE)
	fmt.Println(utils.LEFT_CENTER + strings.Repeat(utils.CENTER, 6) + utils.CENTER_TOP + strings.Repeat(utils.CENTER, specificOptionLength) + utils.RIGHT_CENTER)
	fmt.Println(utils.Parse(specificOptionTpl, struct {
		INDEX  string
		OPTION string
	}{"0.", "Query service mapping port"}))
	fmt.Println(utils.LEFT_CENTER + strings.Repeat(utils.CENTER, 6) + utils.CENTER_CENTER + strings.Repeat(utils.CENTER, maxLength1+11+2) + utils.CENTER_TOP + strings.Repeat(utils.CENTER, maxLength2+8+2) + utils.CENTER_TOP + strings.Repeat(utils.CENTER, maxLength3+2) + utils.RIGHT_CENTER)
	for index, pod := range podList {
		fmt.Println(utils.Parse(optionTpl, struct {
			INDEX  string
			OPTION string
			BRANCH string
			NAME   string
		}{mritd.ToString(index+1) + ".", pod.app, pod.branch, pod.name}))
	}
	fmt.Println(utils.LEFT_CENTER + strings.Repeat(utils.CENTER, 6) + utils.CENTER_CENTER + strings.Repeat(utils.CENTER, maxLength1+11+2) + utils.CENTER_BOTTOM + strings.Repeat(utils.CENTER, maxLength2+8+2) + utils.CENTER_BOTTOM + strings.Repeat(utils.CENTER, maxLength3+2) + utils.RIGHT_CENTER)
	fmt.Println(utils.Parse(specificOptionWithInfoTpl, struct {
		INDEX  string
		OPTION string
		INFO   string
	}{"00.", "Switch Context", "[" + currentContext + "]"}))
	fmt.Println(utils.LEFT_CENTER + strings.Repeat(utils.CENTER, 6) + utils.CENTER_CENTER + strings.Repeat(utils.CENTER, specificOptionLength) + utils.RIGHT_CENTER)
	fmt.Println(utils.Parse(specificOptionWithInfoTpl, struct {
		INDEX  string
		OPTION string
		INFO   string
	}{"000.", "Switch Namespace", "[" + currentNamespace + "]"}))
	fmt.Println(utils.LEFT_BOTTOM + strings.Repeat(utils.CENTER, 6) + utils.CENTER_BOTTOM + strings.Repeat(utils.CENTER, specificOptionLength) + utils.RIGHT_BOTTOM)
	fmt.Println()
	if fastIndex == -1 {
		if maxOption == 1 {
			fmt.Printf(" 请输入数字 [0-%d,00,000](默认-进入这个容器)：", maxOption)
		} else {
			fmt.Printf(" 请输入数字 [0-%d,00,000](默认-退出脚本)：", maxOption)
		}
	}

	var notValid = true
	var input string
	var index int
	for notValid {
		if fastIndex > 0 {
			index = fastIndex
			input = string(index)
			notValid = false
		} else {
			size, err := fmt.Scanln(&input)
			i, err1 := strconv.Atoi(input)
			if isFast && input == "" {
				index = 1
				input = string(index)
				notValid = false
			} else {
				if input == "" {
					if maxOption == 1 {
						input = string(1)
						size = len(input)
						i = 1
					} else {
						fmt.Println()
						fmt.Println(utils.INFO + " 不操作，退出脚本中...")
						os.Exit(0)
					}
				}
				if size < 1 || i < 0 || i > maxOption || err != nil || err1 != nil {
					fmt.Println()
					fmt.Printf(utils.ERROR+" 请输入正确数字 [0-%d,00,000]：", maxOption)
				} else {
					index = i
					notValid = false
				}
			}
		}

	}
	var optionIndex = options(index)
	switch {
	case index == 0 && input == "00": //切换环境
		optionIndex = OptionSwitchContext
	case index == 0 && input == "000": //切换命名空间
		optionIndex = OptionSwitchNamespace
	case index == 0: //查询端口映射
		optionIndex = OptionQueryPort
	default: //进入容器
		optionIndex = OptionEnterPod
	}
	if optionIndex == 0 {
		return optionIndex, podList[index-1]
	} else {
		return optionIndex, Pod{}
	}
}

func selectContainer(containers []Container) Container {
	if containers == nil {
		fmt.Println("没有容器，退出中...")
		os.Exit(0)
	}
	//存储解析后的容器列表
	var containerList []string
	maxLength := 28

	for _, container := range containers {
		name := container.name
		containerList = append(containerList, name)
		maxLength = utils.Max(maxLength, len(name))
	}

	length := maxLength + 17
	cfg := &mritd.SelectConfig{
		ActiveTpl:       utils.LINE + " {{\"↣  Connect to \"|blue}}{{.|lLength " + mritd.ToString(maxLength) + "|magenta}}{{\"  ↢\"|blue}} " + utils.LINE,
		InactiveTpl:     utils.LINE + " {{\"   Connect to \"|white}}{{.|lLength " + mritd.ToString(maxLength) + "|cyan}}    " + utils.LINE,
		SelectedTpl:     utils.INFO_TPL + " 容器名指定为： {{.|cyan}}",
		DisPlaySize:     len(containerList),
		SelectHeaderTpl: utils.LINE + " {{\"Select Container Name:\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		SelectPromptTpl: utils.LINE + " {{\"Use the arrow keys to navigate: ↓ ↑ → ←\"|lLength " + mritd.ToString(length) + "}} " + utils.LINE,
		ShowBorder:      true,
		ShowWidth:       length,
	}
	s := &mritd.Select{
		Items:  containerList,
		Config: cfg,
	}
	return containers[s.Run()]
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
