package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	githubLabel = "github"
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

// RepositoriesFromDeployment returns the reportories attached by 'github' label
func RepositoriesFromDeployment(deployment *v1beta1.Deployment) (map[string]string, error) {
	v, ok := deployment.Labels[githubLabel]
	if !ok {
		return map[string]string{}, errors.Errorf("label %q not found in deployment %q", githubLabel, deployment.Name)
	}

	repos := map[string]string{}

	for _, f := range strings.Split(v, ",") {
		ss := strings.Split(f, ":")
		if len(ss) != 2 {
			return map[string]string{}, errors.Errorf(`invalid label %q value %q, must be "container=owner/repo"`, githubLabel, f)
		}
		repos[ss[0]] = ss[1]
	}

	return repos, nil
}
