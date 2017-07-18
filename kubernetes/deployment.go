package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	githubAnnotation = "github"
)

// Deployment represents the wrapper of Kubernetes Deployment
type Deployment struct {
	raw *v1beta1.Deployment
}

// NewDeployment creates new Deployment object
func NewDeployment(raw *v1beta1.Deployment) *Deployment {
	return &Deployment{
		raw: raw,
	}
}

// Annotations returns the annotations of Deployment
func (d *Deployment) Annotations() map[string]string {
	return d.raw.Annotations
}

// Containers returns the containers inside Deployment
func (d *Deployment) Containers() []*Container {
	containers := []*Container{}

	for i := range d.raw.Spec.Template.Spec.Containers {
		containers = append(containers, NewContainer(&d.raw.Spec.Template.Spec.Containers[i]))
	}

	return containers
}

// Labels returns the labels of Deployment
func (d *Deployment) Labels() map[string]string {
	return d.raw.Labels
}

// Name returns the name of Deployment
func (d *Deployment) Name() string {
	return d.raw.Name
}

// Namespace returns the namespace of Deployment
func (d *Deployment) Namespace() string {
	return d.raw.Namespace
}

// Repositories returns the reportories attached by 'github' annotation
func (d *Deployment) Repositories() (map[string]string, error) {
	v, ok := d.Annotations()[githubAnnotation]
	if !ok {
		return map[string]string{}, errors.Errorf("annotation %q not found in Deployment %q", githubAnnotation, d.Name())
	}

	repos := map[string]string{}

	for _, f := range strings.Split(v, ",") {
		ss := strings.Split(f, "=")
		if len(ss) != 2 {
			return map[string]string{}, errors.Errorf(`invalid annotation %q value %q, must be "container=owner/repo"`, githubAnnotation, f)
		}
		repos[ss[0]] = ss[1]
	}

	return repos, nil
}
