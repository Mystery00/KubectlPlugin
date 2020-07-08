package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

var fastName string
var fastIndexInput string
var fastIndex = -1
var isFast bool

func init() {
	rootCmd.AddCommand(fastCmd)
	rootCmd.PersistentFlags().StringVarP(&fastName, "fast", "f", "", "指定快速进入的服务名称")
	rootCmd.PersistentFlags().StringVarP(&fastIndexInput, "index", "i", "", "指定快速进入的服务序号")
}

var fastCmd = &cobra.Command{
	Use:     "fast",
	Aliases: []string{"f"},
	Short:   "快速进入到服务。",
	Long:    `根据指定的名称快速进入到服务中。`,
	Run: func(cmd *cobra.Command, args []string) {
		isFast = fastName != ""
		i, err := strconv.Atoi(fastIndexInput)
		if err != nil || i <= 0 {
			i = -1
		}
		fastIndex = i
		doAction()
	},
}
