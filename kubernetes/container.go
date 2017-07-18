package kubernetes

import (
	"k8s.io/client-go/pkg/api/v1"
)

// Container represents the wrapper of Kubernetes Pod container
type Container struct {
	raw *v1.Container
}

// NewContainer creates new Container object
func NewContainer(raw *v1.Container) *Container {
	return &Container{
		raw: raw,
	}
}

// Name represents the image name of container
func (c *Container) Image() string {
	return c.raw.Image
}

// Name represents the name of contaienr
func (c *Container) Name() string {
	return c.raw.Name
}
