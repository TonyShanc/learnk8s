// list all-namespaces pods info

package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/etc/rancher/k3s/k3s.yaml")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// all namepsaces pods
	podClient := clientset.CoreV1().Pods("")
	list, err := podClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, pod := range list.Items {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}
}
