// combine dynamic client and discovery client

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	mustInitClient()
}

var (
	discoveryClient *discovery.DiscoveryClient
	dynamicClient   dynamic.Interface
	once            sync.Once

	gvk = &schema.GroupVersionKind{
		Version: "v1",
		Kind:    "Pod",
	}
	namespace = "default"
)

func main() {
	ctx := context.Background()
	resource, err := getDynamicResource(gvk, namespace)
	if err != nil {
		log.Fatal(err)
	}

	create := func() {
		fmt.Println("creating pod......")
		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-pod",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "test-busybox",
						Image: "busybox",
					},
				},
			},
		}
		obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&pod)
		if err != nil {
			log.Fatal(err)
		}
		var unstructured unstructured.Unstructured
		unstructured.Object = obj
		if _, err := resource.Create(ctx, &unstructured, metav1.CreateOptions{}); err != nil {
			log.Fatal(err)
		}
	}
	create()

	get := func() {
		fmt.Println("getting pod......")
		unstructured, err := resource.Get(ctx, "test-pod", metav1.GetOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var pod corev1.Pod
		err = runtime.DefaultUnstructuredConverter.
			FromUnstructured(unstructured.Object, &pod)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(pod.Name)
	}
	get()

	list := func() {
		fmt.Println("listing pods......")
		unstructuredList, err := resource.List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		var podList []*corev1.Pod
		for _, unstructured := range unstructuredList.Items {
			var pod corev1.Pod
			err := runtime.DefaultUnstructuredConverter.
				FromUnstructured(unstructured.Object, &pod)
			if err != nil {
				log.Fatal(err)
			}
			podList = append(podList, &pod)
		}
		for _, pod := range podList {
			fmt.Println(pod.Namespace, pod.Name)
		}
	}
	list()

	delete := func(){
		fmt.Println("deleting pod...")
		if err := resource.Delete(ctx, "test-pod", metav1.DeleteOptions{}); err != nil {
			panic(err)
		}
	}
	delete()
}

func mustInitClient() {
	once.Do(func() {

		var configPath string

		if value := os.Getenv("KUBECONFIG"); value != "" {
			configPath = value
		} else {
			configPath = "~/.kube/config"
		}

		config, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			log.Fatal(err)
		}

		discoveryClient, err = discovery.NewDiscoveryClientForConfig(config)
		if err != nil {
			log.Fatal(err)
		}

		dynamicClient, err = dynamic.NewForConfig(config)
		if err != nil {
			log.Fatal(err)
		}
	})
}

// getDynamicResource convert GVK into dynamic resource
func getDynamicResource(gvk *schema.GroupVersionKind, namespace string) (dr dynamic.ResourceInterface, err error) {

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, fmt.Errorf("CRD has not been registed, err: %s", err)
	}

	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = dynamicClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		dr = dynamicClient.Resource(mapping.Resource)
	}

	return
}
