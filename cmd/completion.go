package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "生成自动补全的脚本",
	Long: `To load completion run

. <(k8s completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(k8s completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = rootCmd.GenZshCompletion(os.Stdout)
	},
}
