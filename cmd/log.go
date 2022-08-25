package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"

	"github.com/spf13/cobra"
)

type LogArgs struct {
	NameSpace string
	Follow    bool
}

var logArgs LogArgs

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [mpijob name]",
	Short: "get job log",
	Long:  `get job log`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := utils.InitKubeClient(); err != nil {
			return
		}
		err := utils.GetNamespace(logArgs.NameSpace)
		if err != nil {
			fmt.Println(err)
			return
		}
		podName := name + "-launcher"
		_ = utils.GetPodLog(podName, logArgs.NameSpace, logArgs.Follow)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVar(&logArgs.NameSpace, "ns", "farctl", "MPI Job Namespace")
	logCmd.Flags().BoolVarP(&logArgs.Follow, "follow", "w", false, "Whether follow the log")
}
