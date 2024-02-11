package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
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
		fmt.Printf("error building kube config: %s\n", err)
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

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return c.Resource(zephySchema).Namespace("").List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return c.Resource(zephySchema).Namespace("").Watch(context.TODO(), options)
			},
		},
		&unstructured.Unstructured{},
		0,
		cache.Indexers{},
	)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			spew.Dump("Resource created: ", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			spew.Dump("Resource: %v was update to: %v", oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			spew.Dump("Resource deleted: %v", obj)
		},
	})

	stop := make(chan struct{})
	defer close(stop)

	go informer.Run(stop)

	if !cache.WaitForCacheSync(stop, informer.HasSynced) {
		panic("Timeout waiting for cache sync")
	}

	fmt.Println("Custom Resource Controller started successfully")

	<-stop
}
