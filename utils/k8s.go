package utils

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	batchv1 "github.com/FFFFFaraway/MPI-Operator/api/batch.test.bdap.com/v1"
	"github.com/FFFFFaraway/MPI-Operator/client/clientset/versioned"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/cmd/cp"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"os"
	"path/filepath"
	"time"
)

var (
	config          *rest.Config
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

	var err error
	// use the current context in kubeconfig
	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	mpijobClientset, err = versioned.NewForConfig(config)
	return err
}

func WaitPodRunning(name, ns string, numWorkers int) {
	fmt.Printf("waiting for %v workers container to be created...\n", numWorkers)
	for {
		time.Sleep(10 * time.Second)
		exit := true
		for i := 0; i < numWorkers; i++ {
			podName := fmt.Sprintf("%v-worker-%v", name, i)
			pod, err := clientset.CoreV1().Pods(ns).Get(context.TODO(), podName, metav1.GetOptions{})
			if err != nil {
				fmt.Println(err)
				return
			}
			if pod.Status.Phase != "Running" {
				fmt.Printf("worker %v is not running\n", podName)
				exit = false
				break
			}
		}
		if exit {
			return
		}
	}
}

// KubectlCp copied from https://stackoverflow.com/questions/51686986/how-to-copy-file-to-container-with-kubernetes-client-go
func KubectlCp(src string, dst string, containername string, ns string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	ioStreams, in, out, errOut := genericclioptions.NewTestIOStreams()
	copyOptions := cp.NewCopyOptions(ioStreams)
	copyOptions.Clientset = clientset
	copyOptions.ClientConfig = config
	copyOptions.Container = containername
	//################### sw edited ########################
	// can't modify namespace by this, it will be overwritten in copyOptions.Complete
	// copyOptions.Namespace = ns
	defaultConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag().WithDiscoveryBurst(300).WithDiscoveryQPS(50.0)
	defaultConfigFlags.Namespace = &ns
	//kubectlOptions := cmd.KubectlOptions{
	//	PluginHandler: cmd.NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes),
	//	Arguments:     os.Args,
	//	ConfigFlags:   defaultConfigFlags,
	//	IOStreams:     genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	//}
	//kubeConfigFlags := kubectlOptions.ConfigFlags
	//if kubeConfigFlags == nil {
	//	kubeConfigFlags = defaultConfigFlags
	//}
	kubeConfigFlags := defaultConfigFlags
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	cmdutil.CheckErr(copyOptions.Complete(f, cp.NewCmdCp(f, ioStreams), []string{src, dst}))
	cmdutil.CheckErr(copyOptions.Validate())
	err := copyOptions.Run()

	//################### sw edited ########################
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not run copy operation: %v", err)
	}
	return in, out, errOut, nil
}
