package cmd

import (
	"fmt"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy image",
	RunE:  doDeploy,
}

var deployOpts = struct {
	container  string
	deployment string
	namespace  string
}{}

func doDeploy(cmd *cobra.Command, args []string) error {
	client, err := kubernetes.NewClient(rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployment, err := client.DetectTargetDeployment(deployOpts.namespace, deployOpts.deployment)
	if err != nil {
		return errors.Wrap(err, "failed to detect target Deployment")
	}

	container, err := client.DetectTargetContainer(deployment, deployOpts.container)
	if err != nil {
		return errors.Wrap(err, "failed to detect target container")
	}

	fmt.Println("deployment: " + deployment.Name)
	fmt.Println("container:  " + container.Name)

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployOpts.container, "container", "c", "", "target container")
	deployCmd.Flags().StringVarP(&deployOpts.deployment, "deployment", "d", "", "target Deployment")
	deployCmd.Flags().StringVar(&deployOpts.namespace, "namespace", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
