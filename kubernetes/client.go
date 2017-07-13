package kubernetes

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client represents the wrapper of Kubernetes API client
type Client struct {
	clientConfig clientcmd.ClientConfig
	clientset    kubernetes.Interface
}

// NewClient creates Client object using local kubecfg
func NewClient(kubeconfig, context string) (*Client, error) {
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
		clientConfig: clientConfig,
		clientset:    clientset,
	}, nil
}

// NewClientInCluster creates Client object in Kubernetes cluster
func NewClientInCluster() (*Client, error) {
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

// DetectTargetDeployment returns the matched or the first deployment
func (c *Client) DetectTargetDeployment(namespace, name string) (*v1beta1.Deployment, error) {
	var deployment *v1beta1.Deployment

	if name == "" {
		ds, err := c.ListDeployment(namespace)
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve Deployments")
		}

		if len(ds) == 0 {
			return nil, errors.Errorf("no Deployment found in namespace %q", namespace)
		}

		if len(ds) > 1 {
			names := []string{}

			for _, d := range ds {
				names = append(names, d.Name)
			}

			return nil, errors.Errorf("multiple Deployments %q found in namespace %q", names, namespace)
		}

		deployment = &ds[0]
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
func (c *Client) GetDeployment(namespace, name string) (*v1beta1.Deployment, error) {
	deployment, err := c.clientset.ExtensionsV1beta1().Deployments(namespace).Get(name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve Deployment %q", name)
	}

	return deployment, nil
}

// ListDeployment returns the list of deployment
func (c *Client) ListDeployment(namespace string) ([]v1beta1.Deployment, error) {
	deployments, err := c.clientset.ExtensionsV1beta1().Deployments(namespace).List(v1.ListOptions{})
	if err != nil {
		return []v1beta1.Deployment{}, errors.Wrap(err, "failed to retrieve Deployments")
	}

	return deployments.Items, nil
}
