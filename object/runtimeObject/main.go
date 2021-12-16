// exercise: transfer runtime.Object and any k8s resource

// runtime.Object is base of k8s type system
// GetObjectKind has been implemented by TypeMeta, we need to focus on DeepCopyObject

package main

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"name": "foo"},
		},
	}

	obj := runtime.Object(pod)

	pod2, ok := obj.(*corev1.Pod)
	if !ok {
		panic("unexpected")
	}
	// check if deep-equal
	if !reflect.DeepEqual(pod, pod2) {
		panic("not equal")
	}
}
