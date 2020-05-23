package cmd

import (
	"KubectlPlugin/mritd"
	"KubectlPlugin/utils"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"strings"
)

func init() {
	rootCmd.AddCommand(portCmd)
}

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "显示节点信息以及服务端口映射列表。",
	Run: func(cmd *cobra.Command, args []string) {
		//快速显示节点信息
		if currentContext == "" {
			//没有设置上下文，并且没有设置命名空间，那么查询一次上下文
			setContextOnly(getCurrentContext())
		}
		if currentNamespace == "" {
			//没有设置命名空间，检测是否有默认命名空间，如果没有弹出菜单指定
			detectNamespace()
		}
		//展示节点信息
		queryPort()
	},
}

type NodeRole string

const (
	MasterString NodeRole = "❆ master ❆"
	NodeString   NodeRole = "-  node  -"
)

func NodeRoleOf(s string) NodeRole {
	switch s {
	case "master":
		return MasterString
	case "node":
		return NodeString
	default:
		return ""
	}
}

type Node struct {
	name string
	role NodeRole
}

type Service struct {
	app       string
	name      string
	clusterIP string
	showPort  string
	ports     []ForwardPort
}

type ForwardPort struct {
	name       string
	nodePort   string
	port       string
	protocol   string
	targetPort string
}

func queryPort() {
	fmt.Println()
	fmt.Println(utils.INFO + " 正在查询节点信息列表...")
	result := utils.Cmd("kubectl", "get", "node", "-o", "json")
	//获取节点列表
	nodes := gjson.Get(result, "items").Array()
	//存储解析后的节点信息
	var nodeList []Node
	for _, node := range nodes {
		name := gjson.Get(node.Raw, "metadata.name").String()
		role := NodeRoleOf(gjson.Get(node.Raw, "metadata.labels.kubernetes\\.io\\/role").String())
		nodeList = append(nodeList, Node{name, role})
	}
	masterTpl := utils.LINE + " [{{.|rLength 15|blue}}] ➜ [{{\"" + string(MasterString) + "\"|lLength 10|magenta}}] " + utils.LINE
	nodeTpl := utils.LINE + " [{{.|rLength 15|blue}}] ➜ [{{\"" + string(NodeString) + "\"|lLength 10|green}}] " + utils.LINE
	fmt.Println(utils.LEFT_TOP + strings.Repeat(utils.CENTER, 34) + utils.RIGHT_TOP)
	for _, node := range nodeList {
		switch node.role {
		case MasterString:
			fmt.Println(utils.Parse(masterTpl, node.name))
		case NodeString:
			fmt.Println(utils.Parse(nodeTpl, node.name))
		}
	}
	fmt.Println(utils.LEFT_BOTTOM + strings.Repeat(utils.CENTER, 34) + utils.RIGHT_BOTTOM)
	fmt.Println()
	fmt.Println(utils.INFO + " 正在查询服务端口映射列表...")
	result = utils.Cmd("kubectl", "get", "service", "-n", currentNamespace, "-o", "json")
	//获取服务列表
	services := gjson.Get(result, "items").Array()
	//存储解析后的服务信息
	var serviceList []Service
	maxLength1 := 3
	maxLength2 := 3
	for _, service := range services {
		app := gjson.Get(service.Raw, "metadata.labels.app").String()
		name := gjson.Get(service.Raw, "metadata.name").String()
		clusterIP := gjson.Get(service.Raw, "spec.clusterIP").String()

		if app == "" {
			app = name + "?"
		}

		var portList []ForwardPort
		for _, port := range gjson.Get(service.Raw, "spec.ports").Array() {
			portName := gjson.Get(port.Raw, "name").String()
			nodePort := gjson.Get(port.Raw, "nodePort").String()
			port1 := gjson.Get(port.Raw, "port").String()
			protocol := gjson.Get(port.Raw, "protocol").String()
			targetPort := gjson.Get(port.Raw, "targetPort").String()
			portList = append(portList, ForwardPort{portName, nodePort, port1, protocol, targetPort})
		}
		var portStringList []string
		for _, port := range portList {
			var portString string
			if port.nodePort == "" { //没有映射端口
				portString = fmt.Sprintf("(%s/%s)", port.protocol, port.port)
			} else {
				portString = fmt.Sprintf("(%s/%s:%s)", port.protocol, port.port, port.nodePort)
			}
			portStringList = append(portStringList, portString)
		}
		showPort := strings.Join(portStringList, ",")
		serviceList = append(serviceList, Service{app, name, clusterIP, showPort, portList})
		maxLength1 = utils.Max(maxLength1, len(app))
		maxLength2 = utils.Max(maxLength2, len(showPort))
	}
	serviceTpl := utils.LINE + " [{{.NAME|rLength " + mritd.ToString(maxLength1) + "|blue}}] ⇛ [{{.PORT|lLength " + mritd.ToString(maxLength2) + "|magenta}}] " + utils.LINE
	length := maxLength1 + maxLength2 + 9
	fmt.Println(utils.LEFT_TOP + strings.Repeat(utils.CENTER, length) + utils.RIGHT_TOP)
	for _, service := range serviceList {
		fmt.Println(utils.Parse(serviceTpl, struct {
			NAME string
			PORT string
		}{service.app, service.showPort}))
	}
	fmt.Println(utils.LEFT_BOTTOM + strings.Repeat(utils.CENTER, length) + utils.RIGHT_BOTTOM)
}
