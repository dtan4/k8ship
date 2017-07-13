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
	namespace string
}{}

func doDeploy(cmd *cobra.Command, args []string) error {
	client, err := kubernetes.NewClient(rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployments, err := client.ListDeployment(deployOpts.namespace)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve deployments")
	}

	for _, d := range deployments {
		fmt.Println(d.Name)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&deployOpts.namespace, "namespace", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
