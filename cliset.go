package k8suite

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	home = homeDir()
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// BuildConfigFromFlags checks envvar `IN_CLUSTER` or `KUBERNETES_SERVICE_HOST` and flags to create `*rest.Config`.
func BuildConfigFromFlags() (*rest.Config, error) {
	var (
		kubeconfig     *string
		currentContext *string
	)
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	currentContext = flag.String("context", "", "kube context")
	flag.Parse()
	return BuildConfig(*kubeconfig, *currentContext)
}

// BuildConfig returns InClusterConfig or LocalConfig automatically.
func BuildConfig(cfgPath, cfgContext string) (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		return rest.InClusterConfig()
	}

	if cfgPath == "" {
		cfgPath = filepath.Join(home, ".kube", "config")
	}

	if cfgContext == "" {
		return clientcmd.BuildConfigFromFlags("", cfgPath)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: cfgPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: cfgContext,
		}).ClientConfig()
}

// NewClientset creates `*kubernetes.Clientset`
func NewClientset() (*kubernetes.Clientset, error) {
	cfg, err := BuildConfigFromFlags()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}

// NewClientsetOrDie returns `*kubernetes.Clientset` and panic on error
func NewClientsetOrDie() *kubernetes.Clientset {
	cliset, err := NewClientset()
	if err != nil {
		panic(err)
	}
	return cliset
}
