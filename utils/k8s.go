package utils

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
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
)

var (
	config    *rest.Config
	clientset *kubernetes.Clientset
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

func GetPod(name, ns string) (*v1.Pod, error) {
	return clientset.CoreV1().Pods(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func GetStatefulSet(name, ns string) (*appsv1.StatefulSet, error) {
	return clientset.AppsV1().StatefulSets(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func PrintPodLog(name, ns string, follow bool) error {
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
	if err != nil {
		panic(err.Error())
	}
	mpijobClientset, err = NewMPIClient()
	return err
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
