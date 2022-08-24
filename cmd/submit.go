package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

type SubmitArgs struct {
	Name              string   `yaml:"name"`
	NameSpace         string   `yaml:"namespace"`
	NumWorkers        int      `yaml:"numWorkers"`
	GitUrl            string   `yaml:"gitUrl"`
	GitRepoName       string   `yaml:"gitRepoName"`
	WorkDir           string   `yaml:"workDir"`
	Commands          []string `yaml:"commands"`
	PipInstall        bool     `yaml:"pipInstall"`
	GpuPerWorker      int      `yaml:"gpuPerWorker"`
	Gang              bool     `yaml:"gang"`
	GangSchedulerName string   `yaml:"gangSchedulerName"`
}

var submitArgs SubmitArgs

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit [mpijob name]",
	Short: "submit a mpi job",
	Long:  `submit a mpi job`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		submitArgs.Name = args[0]
		if submitArgs.NameSpace == "" {
			fmt.Println("namespace needed")
			return
		}
		if submitArgs.GitUrl == "" {
			fmt.Println("git url needed")
			return
		}
		if submitArgs.Gang {
			if submitArgs.GangSchedulerName == "" {
				fmt.Println("use default gang scheduler name: gang-scheduler")
				submitArgs.GangSchedulerName = "gang-scheduler"
			}
		}
		parts := strings.Split(strings.Trim(submitArgs.GitUrl, "/"), "/")
		submitArgs.GitRepoName = strings.Split(parts[len(parts)-1], ".git")[0]
		if len(submitArgs.Commands) == 0 {
			fmt.Println("command needed")
			return
		}
		if err := utils.EnsureNamespace(submitArgs.NameSpace); err != nil {
			return
		}
		if err := utils.InstallRelease(submitArgs.Name, submitArgs.NameSpace, submitArgs); err != nil {
			fmt.Println("helm install error", err)
			return
		}
	},
}

func init() {
	if err := utils.InitKubeClient(); err != nil {
		return
	}
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVar(&submitArgs.NameSpace, "ns", "", "MPI Job Namespace")
	submitCmd.Flags().IntVarP(&submitArgs.NumWorkers, "numWorkers", "n", 1, "Number of Workers")
	submitCmd.Flags().StringVarP(&submitArgs.GitUrl, "gitUrl", "i", "", "git repo link for sync code")
	submitCmd.Flags().StringVar(&submitArgs.WorkDir, "wd", ".", "working directory under project")
	submitCmd.Flags().StringArrayVarP(&submitArgs.Commands, "commands", "c", []string{}, "entry point")
	submitCmd.Flags().BoolVar(&submitArgs.PipInstall, "pip", false, "whether needed to run pip install requirements.txt for workers")
	submitCmd.Flags().IntVar(&submitArgs.GpuPerWorker, "gpu", 1, "number of gpu allocated for each workers")
	submitCmd.Flags().BoolVar(&submitArgs.Gang, "gang", false, "whether use gang scheduler")
	submitCmd.Flags().StringVar(&submitArgs.GangSchedulerName, "gangSchedulerName", "", "gang scheduler name")
}
