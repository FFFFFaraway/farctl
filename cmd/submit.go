package cmd

import (
	"fmt"
	"github.com/FFFFFaraway/farctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

type SubmitArgs struct {
	Name              string   `yaml:"name"`
	Namespace         string   `yaml:"namespace"`
	NumWorkers        int      `yaml:"numWorkers"`
	GitUrl            string   `yaml:"gitUrl"`
	LocalUrl          string   `yaml:"localUrl"`
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
		parts := strings.Split(strings.Trim(submitArgs.GitUrl, "/"), "/")
		submitArgs.GitRepoName = strings.Split(parts[len(parts)-1], ".git")[0]
		if len(submitArgs.Commands) == 0 {
			fmt.Println("command needed")
			return
		}
		if err := utils.InitKubeClient(); err != nil {
			return
		}
		if err := utils.EnsureNamespace(submitArgs.Namespace); err != nil {
			return
		}
		if err := utils.InstallRelease(submitArgs.Name, submitArgs.Namespace, submitArgs); err != nil {
			fmt.Println("helm install error", err)
			return
		}
		fmt.Println("local files uploading...")
		utils.WaitPodRunning(submitArgs.Name, submitArgs.Namespace, submitArgs.NumWorkers)
		for i := 0; i < submitArgs.NumWorkers; i++ {
			podName := fmt.Sprintf("%v-worker-%v", submitArgs.Name, i)
			_, out, _, err := utils.KubectlCp(submitArgs.LocalUrl, podName+":/local-repo", "worker", submitArgs.Namespace)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			fmt.Println("out:")
			fmt.Printf("%s", out.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVar(&submitArgs.Namespace, "ns", "farctl", "MPI Job Namespace")
	submitCmd.Flags().IntVarP(&submitArgs.NumWorkers, "numWorkers", "n", 2, "Number of Workers")

	var url string
	submitCmd.Flags().StringVarP(&url, "codeUrl", "i", ".", "local path or github remote link to sync code")
	if strings.HasPrefix(url, "http") {
		submitArgs.GitUrl = url
	} else {
		submitArgs.LocalUrl = url
	}

	submitCmd.Flags().StringVar(&submitArgs.WorkDir, "wd", ".", "working directory under project")
	submitCmd.Flags().StringArrayVarP(&submitArgs.Commands, "commands", "c", []string{}, "entry point")
	submitCmd.Flags().BoolVar(&submitArgs.PipInstall, "pip", false, "whether needed to run pip install requirements.txt for workers")
	submitCmd.Flags().IntVar(&submitArgs.GpuPerWorker, "gpu", 1, "number of gpu allocated for each workers")
	submitCmd.Flags().BoolVar(&submitArgs.Gang, "gang", false, "whether use gang scheduler")
	submitCmd.Flags().StringVar(&submitArgs.GangSchedulerName, "gangSchedulerName", "gang-scheduler", "gang scheduler name")
}
