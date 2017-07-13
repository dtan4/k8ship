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
	dryRun     bool
	image      string
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

	if deployOpts.dryRun {
		fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", deployment.Name, container.Name)
		fmt.Printf("[dry-run]  before: %s\n", container.Image)
		fmt.Printf("[dry-run]   after: %s\n", deployOpts.image)
	} else {
		fmt.Printf("deploy to (deployment: %q, container: %q)\n", deployment.Name, container.Name)
		fmt.Printf("  before: %s\n", container.Image)
		fmt.Printf("   after: %s\n", deployOpts.image)

		if _, err := client.SetImage(deployment, container.Name, deployOpts.image); err != nil {
			return errors.Wrap(err, "failed to set image")
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&deployOpts.container, "container", "c", "", "target container")
	deployCmd.Flags().StringVarP(&deployOpts.deployment, "deployment", "d", "", "target Deployment")
	deployCmd.Flags().BoolVar(&deployOpts.dryRun, "dry-run", false, "dry run")
	deployCmd.Flags().StringVar(&deployOpts.image, "image", "", "new image")
	deployCmd.Flags().StringVar(&deployOpts.namespace, "namespace", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
