package utils

import (
	"context"
	"fmt"
	batchv1 "github.com/FFFFFaraway/MPI-Operator/api/batch.test.bdap.com/v1"
	"github.com/FFFFFaraway/MPI-Operator/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"text/tabwriter"
	"time"
)

var mpijobClientset *versioned.Clientset

func GetMPIJob(name, ns string) (*batchv1.MPIJob, error) {
	return mpijobClientset.BatchV1().MPIJobs(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func ListMPIJob(ns string) (*batchv1.MPIJobList, error) {
	return mpijobClientset.BatchV1().MPIJobs(ns).List(context.TODO(), metav1.ListOptions{})
}

func PrintMPIJobLog(mpijob *batchv1.MPIJob, follow bool) {
	podName := mpijob.Name + "-launcher"
	_ = PrintPodLog(podName, mpijob.Namespace, follow)
}

func WaitMPIJobPodRunning(mpijob *batchv1.MPIJob) error {
	name, ns, numWorkers := mpijob.Name, mpijob.Namespace, int(*mpijob.Spec.NumWorkers)
	fmt.Printf("waiting for %v workers to be created...\n", numWorkers)
	return wait.PollImmediate(10*time.Second, 10*time.Minute, func() (done bool, err error) {
		done = true
		for i := 0; i < numWorkers; i++ {
			podName := fmt.Sprintf("%v-worker-%v", name, i)
			pod, err := GetPod(podName, ns)
			if err != nil {
				return true, err
			}
			if pod.Status.Phase != "Running" {
				fmt.Printf("worker %v is not running\n", podName)
				done = false
				break
			}
		}
		return done, nil
	})
}

func CopyLocalRepoToMPIJob(url string, mpijob *batchv1.MPIJob) {
	fmt.Println("local files uploading...")
	for i := 0; i < int(*mpijob.Spec.NumWorkers); i++ {
		podName := fmt.Sprintf("%v-worker-%v", mpijob.Name, i)
		_, _, _, err := KubectlCp(url, podName+":/local-repo", "worker", mpijob.Namespace)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func PrintMPIJobList(mpijobs *batchv1.MPIJobList) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "Namespace\tName\tReadyWorkers/Total\tLauncherStatus\tAge")
	for _, mpijob := range mpijobs.Items {
		ss, err := GetStatefulSet(mpijob.Name+"-worker", mpijob.Namespace)
		if err != nil {
			fmt.Println(err)
			return
		}
		var launcherStatus string
		l, err := GetPod(mpijob.Name+"-launcher", mpijob.Namespace)
		if err != nil {
			if !errors.IsNotFound(err) {
				fmt.Println(err)
				return
			}
			launcherStatus = "WaitingWorkers"
		} else {
			launcherStatus = string(l.Status.Phase)
		}
		_, _ = fmt.Fprintf(w, "%v\t%v\t%v/%v\t%v\t%v\n",
			mpijob.Namespace,
			mpijob.Name,
			ss.Status.ReadyReplicas,
			ss.Status.Replicas,
			launcherStatus,
			time.Since(mpijob.CreationTimestamp.Time).Round(time.Second),
		)
	}
	_ = w.Flush()
}

func NewMPIClient() (*versioned.Clientset, error) {
	return versioned.NewForConfig(config)
}
