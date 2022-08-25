package utils

import (
	"fmt"
	batchv1 "github.com/FFFFFaraway/MPI-Operator/api/batch.test.bdap.com/v1"
	"os"
	"text/tabwriter"
)

func PrintMPIJobHeader(w *tabwriter.Writer) {
	fmt.Fprintln(w, "Name\tNamespace\tCreationTimestamp\tNumWorkers")
}

func PrintMPIJob(mpijob *batchv1.MPIJob) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	PrintMPIJobHeader(w)
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", mpijob.Name, mpijob.Namespace, mpijob.CreationTimestamp, *mpijob.Spec.NumWorkers)
	w.Flush()
}

func PrintMPIJobList(mpijobs *batchv1.MPIJobList) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	PrintMPIJobHeader(w)
	for _, mpijob := range mpijobs.Items {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", mpijob.Name, mpijob.Namespace, mpijob.CreationTimestamp, *mpijob.Spec.NumWorkers)
	}
	w.Flush()
}
