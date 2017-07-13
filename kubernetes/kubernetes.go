package kubernetes

import (
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

// DefaultConfigFile returns the default kubeconfig file path
func DefaultConfigFile() string {
	return clientcmd.RecommendedHomeFile
}

// DefaultNamespace returns the default namespace
func DefaultNamespace() string {
	return v1.NamespaceDefault
}

// ContainerImageFromDeployment returns image name of the given container
func ContainerImageFromDeployment(deployment *v1beta1.Deployment, container string) string {
	for _, c := range deployment.Spec.Template.Spec.Containers {
		if c.Name == container {
			return c.Image
		}
	}

	return ""
}
