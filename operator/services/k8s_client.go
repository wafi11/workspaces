package services

import (
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	Client    kubernetes.Interface
	DynClient dynamic.Interface
	Mapper    meta.RESTMapper
}

func NewK8sClient(kubeConfigPath string) (*K8sClient, error) {
	log.Printf("loading kubeconfig from: %s", kubeConfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build k8s config: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}

	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

	return &K8sClient{
		Client:    client,
		DynClient: dynClient,
		Mapper:    mapper,
	}, nil
}
