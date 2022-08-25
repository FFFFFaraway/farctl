package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
)

type DeleteArgs struct {
	Namespace string
}

var deleteArgs DeleteArgs

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [mpijob name]",
	Short: "delete a mpi job",
	Long:  `delete a mpi job`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.InitKubeClient(); err != nil {
			return
		}
		err := utils.GetNamespace(deleteArgs.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := utils.DeleteRelease(args[0], deleteArgs.Namespace); err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteArgs.Namespace, "ns", "farctl", "MPI Job Namespace")
}
