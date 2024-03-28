### ⚡ To get up and running 

##### Create the 'Kind' Cluster
```bash
kind create cluster --name my-cluster
```
##### Or provide a config for the cluster
```bash
kind create cluster --name my-cluster --config kind-config.yaml
```

##### Apply the CustomResourceDefinition
```bash
kubectl apply -f zephy-resource-definition.yaml
```

##### Applying the custom resource
```bash
kubectl apply -f zephy-resource.yaml
```

I wanted to understand how kubernetes work and since I am still in my process of learning about it, I thought that trying to create something simple like a CRD would give me a better understaning of the underlying architecture of the components of kubernetes.

All credits go to: https://medium.com/@disha.20.10/building-and-extending-kubernetes-a-writing-first-custom-controller-with-go-bc57a50d61f7
