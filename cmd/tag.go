package cmd

import (
	"fmt"
	"strings"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Deploy by image tagn",
	RunE:  doTag,
}

var tagOpts = struct {
	container  string
	deployment string
	dryRun     bool
	namespace  string
}{}

func doTag(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("image tag must be given")
	}
	tag := args[0]

	client, err := kubernetes.NewClient(rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployment, err := client.DetectTargetDeployment(tagOpts.namespace, tagOpts.deployment)
	if err != nil {
		return errors.Wrap(err, "failed to detect target Deployment")
	}

	container, err := client.DetectTargetContainer(deployment, tagOpts.container)
	if err != nil {
		return errors.Wrap(err, "failed to detect target container")
	}

	currentImage := deployment.ContainerImage(container.Name())
	newImage := strings.Split(currentImage, ":")[0] + ":" + tag

	if tagOpts.dryRun {
		fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("[dry-run]  before: %s\n", container.Image())
		fmt.Printf("[dry-run]   after: %s\n", newImage)
	} else {
		fmt.Printf("deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("  before: %s\n", container.Image())
		fmt.Printf("   after: %s\n", newImage)

		if _, err := client.SetImage(
			deployment, container.Name(), newImage, composeTagCause(tag, container.Name(), deployment.Name(), tagOpts.namespace),
		); err != nil {
			return errors.Wrap(err, "failed to set image")
		}
	}

	return nil
}

func composeTagCause(tag, container, deployment, namespace string) string {
	return fmt.Sprintf(`k8ship tag %s --container "%s" --deployment "%s" --namespace "%s"`, tag, container, deployment, namespace)
}

func init() {
	RootCmd.AddCommand(tagCmd)

	tagCmd.Flags().StringVarP(&tagOpts.container, "container", "c", "", "target container")
	tagCmd.Flags().StringVarP(&tagOpts.deployment, "deployment", "d", "", "target Deployment")
	tagCmd.Flags().BoolVar(&tagOpts.dryRun, "dry-run", false, "dry run")
	tagCmd.Flags().StringVarP(&tagOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
