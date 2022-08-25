package utils

import (
	"bufio"
	"context"
	"flag"
	batchv1 "github.com/FFFFFaraway/MPI-Operator/api/batch.test.bdap.com/v1"
	"github.com/FFFFFaraway/MPI-Operator/client/clientset/versioned"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var (
	mpijobClientset *versioned.Clientset
	clientset       *kubernetes.Clientset
)

func createNamespace(namespace string) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

func GetNamespace(ns string) error {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), ns, metav1.GetOptions{})
	return err
}

func EnsureNamespace(ns string) error {
	err := GetNamespace(ns)
	if err != nil && errors.IsNotFound(err) {
		if err = createNamespace(ns); err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func GetMPIJob(name, ns string) (*batchv1.MPIJob, error) {
	return mpijobClientset.BatchV1().MPIJobs(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func ListMPIJob(ns string) (*batchv1.MPIJobList, error) {
	return mpijobClientset.BatchV1().MPIJobs(ns).List(context.TODO(), metav1.ListOptions{})
}

func GetPodLog(name, ns string, follow bool) error {
	req := clientset.CoreV1().Pods(ns).GetLogs(name, &v1.PodLogOptions{
		Follow: follow,
	})
	r, err := req.Stream(context.TODO())
	if err != nil {
		return err
	}
	buf := bufio.NewWriter(os.Stdout)
	_, err = io.Copy(buf, r)
	if err != nil {
		return err
	}
	return nil
}

func InitKubeClient() error {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	mpijobClientset, err = versioned.NewForConfig(config)
	return err
}
