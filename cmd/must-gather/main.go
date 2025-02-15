// must-gather/cmd/must-gather/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/must-gather/pkg/gather"
	corev1 "k8s.io/api/core/v1"  // Now properly used
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var (
		node        = flag.String("node", "", "Node name to schedule the pod")
		image       = flag.String("image", "quay.io/openshift/must-gather:latest", "Container image")
		hasMaster   = flag.Bool("has-master", false, "Target master nodes")
		sourceDir   = flag.String("source-dir", "/must-gather", "Data directory")
		volumePct   = flag.Int("volume-pct", 80, "Max volume usage percentage")
		hostNetwork = flag.Bool("host-network", false, "Use host networking")
		since       = flag.Duration("since", 0, "Logs newer than duration")
		sinceTime   = flag.String("since-time", "", "Logs after timestamp")
		namespace   = flag.String("namespace", "default", "Namespace for the pod")
	)
	flag.Parse()

	// Get Kubernetes config
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if envcfg := os.Getenv("KUBECONFIG"); envcfg != "" {
		kubeconfig = envcfg
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	options := &gather.MustGatherOptions{
		Command:          flag.Args(),
		SourceDir:        *sourceDir,
		VolumePercentage: *volumePct,
		HostNetwork:      *hostNetwork,
		Since:            *since,
		SinceTime:        *sinceTime,
	}

	// Explicit type declaration to use corev1.Pod
	var pod *corev1.Pod = options.NewPod(*node, *image, *hasMaster)

	// Create the pod in the cluster
	createdPod, err := clientset.CoreV1().Pods(*namespace).Create(
		context.TODO(),
		pod,
		metav1.CreateOptions{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create pod: %v", err))
	}

	fmt.Printf("Successfully created pod %q in namespace %q\n", createdPod.Name, *namespace)
}