package cmd

import (
	"fmt"
	"strings"

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
	tag        string
}{}

func doDeploy(cmd *cobra.Command, args []string) error {
	if deployOpts.image != "" && deployOpts.tag != "" {
		return errors.New("both target image and tag cannot be specified simultaneously")
	}

	if deployOpts.image == "" && deployOpts.tag == "" {
		return errors.New("target image (--image) or tag (--tag) must be specified")
	}

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

	var newImage string

	if deployOpts.image == "" && deployOpts.tag != "" {
		currentImage := kubernetes.ContainerImageFromDeployment(deployment, container.Name)
		newImage = strings.Split(currentImage, ":")[0] + ":" + deployOpts.tag
	} else {
		newImage = deployOpts.image
	}

	if deployOpts.dryRun {
		fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", deployment.Name, container.Name)
		fmt.Printf("[dry-run]  before: %s\n", container.Image)
		fmt.Printf("[dry-run]   after: %s\n", newImage)
	} else {
		fmt.Printf("deploy to (deployment: %q, container: %q)\n", deployment.Name, container.Name)
		fmt.Printf("  before: %s\n", container.Image)
		fmt.Printf("   after: %s\n", newImage)

		if _, err := client.SetImage(deployment, container.Name, newImage); err != nil {
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
	deployCmd.Flags().StringVar(&deployOpts.tag, "tag", "", "new tag")
}
