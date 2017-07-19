package cmd

import (
	"fmt"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Deploy by image",
	RunE:  doImage,
}

var imageOpts = struct {
	container  string
	deployment string
	dryRun     bool
	namespace  string
}{}

func doImage(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("image must be given")
	}
	image := args[0]

	client, err := kubernetes.NewClient(rootOpts.annotationPrefix, rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployment, err := client.DetectTargetDeployment(imageOpts.namespace, imageOpts.deployment)
	if err != nil {
		return errors.Wrap(err, "failed to detect target Deployment")
	}

	container, err := client.DetectTargetContainer(deployment, imageOpts.container)
	if err != nil {
		return errors.Wrap(err, "failed to detect target container")
	}

	if imageOpts.dryRun {
		fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("[dry-run]  before: %s\n", container.Image())
		fmt.Printf("[dry-run]   after: %s\n", image)
	} else {
		fmt.Printf("deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("  before: %s\n", container.Image())
		fmt.Printf("   after: %s\n", image)

		if _, err := client.SetImage(
			deployment, container.Name(), image, composeImageCause(image, container.Name(), deployment.Name(), tagOpts.namespace),
		); err != nil {
			return errors.Wrap(err, "failed to set image")
		}

		fmt.Printf("\n")
		fmt.Printf("deployment successfully updated! check rollout status by `kubectl rollout status deployment/DEPLOYMENT --namespace %s`\n", imageOpts.namespace)
	}

	return nil
}

func composeImageCause(image, container, deployment, namespace string) string {
	return fmt.Sprintf(`k8ship image %s --container "%s" --deployment "%s" --namespace "%s"`, image, container, deployment, namespace)
}

func init() {
	RootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringVarP(&imageOpts.container, "container", "c", "", "target container")
	imageCmd.Flags().StringVarP(&imageOpts.deployment, "deployment", "d", "", "target Deployment")
	imageCmd.Flags().BoolVar(&imageOpts.dryRun, "dry-run", false, "dry run")
	imageCmd.Flags().StringVarP(&imageOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
