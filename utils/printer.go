package utils

import (
	"context"
	"fmt"
	batchv1 "github.com/FFFFFaraway/MPI-Operator/api/batch.test.bdap.com/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"text/tabwriter"
)

func PrintMPIJobList(mpijobs *batchv1.MPIJobList) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "Namespace\tName\tReadyWorkers/Total\tLauncherStatus\tCreationTimestamp")
	for _, mpijob := range mpijobs.Items {
		ss, err := clientset.AppsV1().StatefulSets(mpijob.Namespace).Get(context.TODO(), mpijob.Name+"-worker", metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		var launcherStatus string
		l, err := clientset.CoreV1().Pods(mpijob.Namespace).Get(context.TODO(), mpijob.Name+"-launcher", metav1.GetOptions{})
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
			mpijob.CreationTimestamp)
	}
	_ = w.Flush()
}
