package kubernetes

import (
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	deployTargetAnnotation          = "deploy-target"
	deployTargetContainerAnnotation = "deploy-target-container"

	githubAnnotation = "github"
)

var (
	deployTargetAnnotationTrue = []string{"1", "true"}
)

// Deployment represents the wrapper of Kubernetes Deployment
type Deployment struct {
	annotationPrefix string
	raw              *v1beta1.Deployment
}

// NewDeployment creates new Deployment object
func NewDeployment(annotationPrefix string, raw *v1beta1.Deployment) *Deployment {
	return &Deployment{
		annotationPrefix: annotationPrefix,
		raw:              raw,
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

// ContainerImage returns image name of the given container
func (d *Deployment) ContainerImage(container string) string {
	for _, c := range d.Containers() {
		if c.Name() == container {
			return c.Image()
		}
	}

	return ""
}

// DeployTargetContainer returns
// - specified in `deploy-target-container` annotation
func (d *Deployment) DeployTargetContainer() (*Container, error) {
	v, ok := d.Annotations()[d.annotationPrefix+deployTargetContainerAnnotation]
	if !ok {
		return nil, errors.Errorf(`annotation "deploy-target-container" does not exist in Deployment %q`, d.Name())
	}

	for _, c := range d.Containers() {
		if c.Name() == v {
			return c, nil
		}
	}

	return nil, errors.Errorf("container %q does not exist in Deployment %q", v, d.Name())
}

// IsDeployTarget returns whether this deployment is deploy target or not
// - has `deploy-target: 1` or `deploy-target: true` annotation
func (d *Deployment) IsDeployTarget() bool {
	for _, v := range deployTargetAnnotationTrue {
		if d.raw.Annotations[d.annotationPrefix+deployTargetAnnotation] == v {
			return true
		}
	}

	return false
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
	v, ok := d.Annotations()[d.annotationPrefix+githubAnnotation]
	if !ok {
		return map[string]string{}, errors.Errorf("annotation %q not found in Deployment %q", d.annotationPrefix+githubAnnotation, d.Name())
	}

	repos := map[string]string{}

	for _, f := range strings.Split(v, ",") {
		ss := strings.Split(f, "=")
		if len(ss) != 2 {
			return map[string]string{}, errors.Errorf(`invalid annotation %q value %q, must be "container=owner/repo"`, d.annotationPrefix+githubAnnotation, f)
		}
		repos[ss[0]] = ss[1]
	}

	return repos, nil
}
