package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"

	"github.com/spf13/cobra"
)

type LogArgs struct {
	Namespace string
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
		err := utils.GetNamespace(logArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		mpijob, err := utils.GetMPIJob(name, submitArgs.Namespace)
		if err != nil {
			fmt.Println("get mpijob error:", err)
			return
		}
		utils.PrintMPIJobLog(mpijob, logArgs.Follow)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().StringVar(&logArgs.Namespace, "ns", "farctl", "MPI Job Namespace")
	logCmd.Flags().BoolVarP(&logArgs.Follow, "follow", "w", true, "Whether follow the log")
}
