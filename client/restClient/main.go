// list all-namespaces pods info

package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/etc/rancher/k3s/k3s.yaml")
	if err != nil {
		panic(err)
	}

	config.APIPath = "api"
	// group core version v1
	config.GroupVersion = &corev1.SchemeGroupVersion
	// set codec
	config.NegotiatedSerializer = scheme.Codecs
	
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	result := &corev1.PodList{}
	
	ctx := context.Background()

	err = restClient.Get().
		Namespace("default").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{Limit: 10}, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	if err != nil {
		panic(err)
	}

	for _, pod := range result.Items {
		fmt.Println(pod.Namespace, pod.Name, pod.Status.Phase)
	}
}