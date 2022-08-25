package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
)

type GetArgs struct {
	Namespace string
}

var getArgs GetArgs

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [mpijob name]",
	Short: "get mpijob configuration",
	Long:  `get mpijob configuration`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.InitKubeClient(); err != nil {
			return
		}
		err := utils.GetNamespace(getArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = utils.GetReleaseValue(args[0], getArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getArgs.Namespace, "ns", "farctl", "MPI Job Namespace")
}
