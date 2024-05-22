package main

import (
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubebuilder-declarative-pattern/mockkubeapiserver"
)

func main() {
	var wg sync.WaitGroup

	k8s, err := mockkubeapiserver.NewMockKubeAPIServer(":55441")
	if err != nil {
		log.Fatalf("error creating mock-apiserver: %v", err)
	}

	k8s.RegisterType(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}, "namespaces", meta.RESTScopeRoot)
	k8s.RegisterType(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}, "secrets", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}, "configmaps", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1alpha2", Kind: "Pod"}, "pods", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Node"}, "nodes", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "foozer.example.com", Version: "v1alpha1", Kind: "FoozerConfig"}, "foozerconfigs", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1alpha2", Kind: "ResourceClass"}, "resourceclasses", meta.RESTScopeRoot)
	k8s.RegisterType(schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1alpha2", Kind: "ResourceClaim"}, "resourceclaims", meta.RESTScopeNamespace)
	k8s.RegisterType(schema.GroupVersionKind{Group: "resource.k8s.io", Version: "v1alpha1", Kind: "ResourcePool"}, "resourcepools", meta.RESTScopeRoot)

	wg.Add(1)
	addr, err := k8s.StartServing()
	if err != nil {
		log.Fatalf("error starting mock-apiserver: %v", err)
	}
	log.Println("addr = ", addr)

	wg.Wait()
}
