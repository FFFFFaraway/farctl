package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
)

var namespace string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [mpijob name]",
	Short: "delete a mpi job",
	Long:  `delete a mpi job`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
			fmt.Println("namespace needed")
			return
		}
		err := utils.GetNamespace(namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := utils.DeleteRelease(args[0], namespace); err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&namespace, "ns", "", "MPI Job Namespace")
}
