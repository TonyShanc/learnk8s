// list all-namespaces pods info

package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/etc/rancher/k3s/k3s.yaml")
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	// all namespaces
	unstructuredObjs, err := dynamicClient.Resource(gvr).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	podList := &corev1.PodList{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObjs.UnstructuredContent(), podList); err != nil {
		panic(err)
	}
	
	for _, pod := range podList.Items {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}
}
