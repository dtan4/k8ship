package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	deployTargetAnnotation          = "deploy-target"
	deployTargetContainerAnnotation = "deploy-target-container"
	deployUserAnnotation            = "deploy-user"
	reloadedAtAnnotation            = "reloaded-at"

	changeCauseAnnotation = "kubernetes.io/change-cause"
	revisionAnnotation    = "deployment.kubernetes.io/revision"

	githubAnnotation = "github"
)

// DefaultConfigFile returns the default kubeconfig file path
func DefaultConfigFile() string {
	return clientcmd.RecommendedHomeFile
}

// DefaultNamespace returns the default namespace
func DefaultNamespace() string {
	return v1.NamespaceDefault
}

// GetTargetImage returns the unique image name of target containers
func GetTargetImage(containers map[string]*Container) (string, error) {
	images := map[string]bool{}

	for _, c := range containers {
		images[strings.Split(c.Image(), ":")[0]] = true
	}

	ss := make([]string, 0, len(images))

	for k := range images {
		ss = append(ss, k)
	}

	if len(images) == 0 {
		return "", errors.Errorf("no image found")
	}

	if len(images) > 1 {
		return "", errors.Errorf("multiple images %q found, all target containers must use the same image", ss)
	}

	return ss[0], nil
}

// GetTargetRepository returns the unique GitHub repository of target containers
func GetTargetRepository(deployments []*Deployment, containers map[string]*Container) (string, error) {
	repos := map[string]bool{}

	for _, d := range deployments {
		c, ok := containers[d.Name()]
		if !ok {
			return "", errors.Errorf("no container found in Deployment %q", d.Name())
		}

		rs, err := d.Repositories()
		if err != nil {
			return "", errors.Wrapf(err, "failed to retrieve repositories of Deployment %q", d.Name())
		}

		v, ok := rs[c.Name()]
		if !ok {
			return "", errors.Errorf("GitHub repository for container %q in Deployment %q is not set", c.Name(), d.Name())
		}
		repos[v] = true
	}

	ss := make([]string, 0, len(repos))

	for k := range repos {
		ss = append(ss, k)
	}

	if len(repos) > 1 {
		return "", errors.Errorf("multiple repositories %q found, all target containers must use the same repository", ss)
	}

	return ss[0], nil
}
