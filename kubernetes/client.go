package kubernetes

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client represents the wrapper of Kubernetes API client
type Client struct {
	annotationPrefix string
	clientConfig     clientcmd.ClientConfig
	clientset        kubernetes.Interface
}

// NewClient creates Client object using local kubecfg
func NewClient(annotationPrefix string, kubeconfig, context string) (*Client, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: context})

	config, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "falied to load local kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load clientset")
	}

	return &Client{
		annotationPrefix: annotationPrefix,
		clientConfig:     clientConfig,
		clientset:        clientset,
	}, nil
}

// NewClientInCluster creates Client object in Kubernetes cluster
func NewClientInCluster(annotationPrefix string) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubeconfig in cluster")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "falied to load clientset")
	}

	return &Client{
		clientset: clientset,
	}, nil
}

// DetectTargetContainer returns the matched or the first container
func (c *Client) DetectTargetContainer(deployment *Deployment, name string) (*Container, error) {
	if name == "" {
		if len(deployment.Containers()) > 1 {
			names := []string{}

			for _, c := range deployment.Containers() {
				names = append(names, c.Name())
			}

			return nil, errors.Errorf("multiple containers %q found in deployment %q", names, deployment.Name())
		}

		return deployment.Containers()[0], nil
	}

	for _, c := range deployment.Containers() {
		if c.Name() == name {
			return c, nil
		}
	}

	return nil, errors.Errorf("container %q does not exist", name)
}

// DetectTargetDeployment returns the matched or the first deployment
func (c *Client) DetectTargetDeployment(namespace, name string) (*Deployment, error) {
	var deployment *Deployment

	if name == "" {
		ds, err := c.ListDeployments(namespace)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve Deployments")
		}

		if len(ds) == 0 {
			return nil, errors.Errorf("no Deployment found in namespace %q", namespace)
		}

		if len(ds) > 1 {
			names := []string{}

			for _, d := range ds {
				names = append(names, d.Name())
			}

			return nil, errors.Errorf("multiple Deployments %q found in namespace %q", names, namespace)
		}

		deployment = ds[0]
	} else {
		d, err := c.GetDeployment(namespace, name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve Deployment %q", name)
		}

		deployment = d
	}

	return deployment, nil
}

// GetDeployment returns a deployment
func (c *Client) GetDeployment(namespace, name string) (*Deployment, error) {
	deployment, err := c.clientset.ExtensionsV1beta1().Deployments(namespace).Get(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve Deployment %q", name)
	}

	return NewDeployment(c.annotationPrefix, deployment), nil
}

// ListDeployments returns the list of deployment
func (c *Client) ListDeployments(namespace string) ([]*Deployment, error) {
	deployments, err := c.clientset.ExtensionsV1beta1().Deployments(namespace).List(v1.ListOptions{})
	if err != nil {
		return []*Deployment{}, errors.Wrap(err, "failed to retrieve Deployments")
	}

	ds := []*Deployment{}

	// `for _, d := range deployments` uses the same pointer in `d`
	for i := range deployments.Items {
		ds = append(ds, NewDeployment(c.annotationPrefix, &deployments.Items[i]))
	}

	return ds, nil
}

// ListReplicaSets returns the list of ReplicaSets
func (c *Client) ListReplicaSets(deployment *Deployment) ([]*ReplicaSet, error) {
	all, err := c.clientset.ExtensionsV1beta1().ReplicaSets(deployment.Namespace()).List(v1.ListOptions{})
	if err != nil {
		return []*ReplicaSet{}, errors.Wrapf(err, "failed to retrieve ReplicaSets")
	}

	filtered := make([]*ReplicaSet, 0, len(all.Items))

	for _, rs := range all.Items {
		for _, or := range rs.GetOwnerReferences() {
			if string(or.UID) == deployment.UID() {
				r := rs
				filtered = append(filtered, NewReplicaSet(&r))
			}
		}
	}

	return filtered, nil
}

// ReloadPods reloads all Pods in the given deployment by setting new annotation
func (c *Client) ReloadPods(deployment *Deployment, signature string) (*Deployment, error) {
	patch := fmt.Sprintf(`{
  "spec": {
    "template": {
      "metadata": {
        "annotations": {
          "%s": %q
        }
      }
    }
  }
}`, c.annotationPrefix+"reloaded-at", signature)

	newd, err := c.clientset.ExtensionsV1beta1().Deployments(deployment.Namespace()).Patch(deployment.Name(), api.StrategicMergePatchType, []byte(patch))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update deployment %q", deployment.Name())
	}

	return NewDeployment(c.annotationPrefix, newd), nil
}

// SetImage sets new image to the given deployments
func (c *Client) SetImage(deployment *Deployment, container, image, user, cause string) (*Deployment, error) {
	patch := fmt.Sprintf(`{
  "metadata": {
    "annotations": {
      "%s": %q
    }
  },
  "spec": {
    "template": {
      "metadata": {
        "annotations": {
          "%s": %q
        }
      },
      "spec": {
        "containers": [
          {
            "name": "%s",
            "image": "%s"
          }
        ]
      }
    }
  }
}`, changeCauseAnnotation, cause, c.annotationPrefix+deployUserAnnotation, user, container, image)

	newd, err := c.clientset.ExtensionsV1beta1().Deployments(deployment.Namespace()).Patch(deployment.Name(), api.StrategicMergePatchType, []byte(patch))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update deployment %q", deployment.Name())
	}

	return NewDeployment(c.annotationPrefix, newd), nil
}
