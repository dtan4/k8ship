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

func doReload(cmd *cobra.Command, args []string) error {
	k8sClient, err := kubernetes.NewClient(rootOpts.annotationPrefix, rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployments, err := k8sClient.ListDeployments(deployOpts.namespace)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Deployments")
	}

	if len(deployments) == 0 {
		return errors.Errorf("no Deployment found in namespace %s", deployOpts.namespace)
	}

	targetDeployments := []*kubernetes.Deployment{}

	for _, d := range deployments {
		if d.IsDeployTarget() {
			targetDeployments = append(targetDeployments, d)
		}
	}

	if len(targetDeployments) == 0 {
		return errors.New("no target Deployments found")
	}

	timestamp := time.Now().Local().String()

	for _, d := range targetDeployments {
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
}
