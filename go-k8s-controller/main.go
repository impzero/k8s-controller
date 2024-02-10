package main

import (
	"context"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var path string
	if home := homedir.HomeDir(); home != "" {
		path = filepath.Join(home, "/.kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		fmt.Errorf("Error building kube config: %s", err)

		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	c, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	zephySchema := schema.GroupVersionResource{
		Group:    "myk8s.io",
		Version:  "v1",
		Resource: "zephy",
	}

	informer := cache.NewSharedIndexInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return c.Resource(zephySchema).List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.Resource(zephySchema).Watch(context.TODO(), options)
		},
		DisableChunking: false,
	}, &metav1.Pod{}, 0, cache.Indexers{})
}
