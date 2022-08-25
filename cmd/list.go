package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
)

type ListArgs struct {
	Namespace string
}

var listArgs ListArgs

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [mpijob namespace]",
	Short: "list all mpijob in namespace",
	Long:  `list all mpijob in namespace`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.InitKubeClient(); err != nil {
			return
		}
		err := utils.GetNamespace(listArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}

		mpiJobs, err := utils.ListMPIJob(listArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		utils.PrintMPIJobList(mpiJobs)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listArgs.Namespace, "ns", "farctl", "MPI Job Namespace")
}
