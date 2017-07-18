package kubernetes

import (
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
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
