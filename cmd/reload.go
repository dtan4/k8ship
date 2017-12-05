package cmd

import (
	"fmt"
	"time"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload Pods in Deployment",
	RunE:  doReload,
}

var reloadOpts = struct {
	deployment string
	namespace  string
}{}

func doReload(cmd *cobra.Command, args []string) error {
	k8sClient, err := kubernetes.NewClient(rootOpts.annotationPrefix, rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	var deployments []*kubernetes.Deployment

	if reloadOpts.deployment == "" {
		ds, err := k8sClient.ListDeployments(reloadOpts.namespace)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve Deployments")
		}

		if len(ds) == 0 {
			return errors.Errorf("no Deployment found in namespace %s", reloadOpts.namespace)
		}

		deployments = ds
	} else {
		d, err := k8sClient.GetDeployment(reloadOpts.namespace, reloadOpts.deployment)
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve Deployment %s in %s", reloadOpts.deployment, reloadOpts.namespace)
		}

		deployments = []*kubernetes.Deployment{d}
	}

	timestamp := time.Now().Local().String()

	for _, d := range deployments {
		_, err := k8sClient.ReloadPods(d, timestamp)
		if err != nil {
			return errors.Wrap(err, "failed to set annotations")
		}

		fmt.Printf("reloaded all Pods in %s\n", d.Name())
	}

	return nil
}

func init() {
	RootCmd.AddCommand(reloadCmd)

	reloadCmd.Flags().StringVarP(&reloadOpts.deployment, "deployment", "d", "", "target Deployment")
	reloadCmd.Flags().StringVarP(&reloadOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
