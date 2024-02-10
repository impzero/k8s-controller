package main

import (
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
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

	_ = schema.GroupVersionResource{
		Group:    "myk8s.io",
		Version:  "v1",
		Resource: "zephy",
	}
}
