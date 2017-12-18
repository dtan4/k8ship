package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dtan4/k8ship/github"
	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// refCmd represents the ref command
var refCmd = &cobra.Command{
	Use:   "ref BRANCH|COMMIT_SHA1",
	Short: "Deploy by git ref",
	RunE:  doRef,
}

var refOpts = struct {
	accessToken string
	container   string
	deployment  string
	dryRun      bool
	namespace   string
	user        string
}{}

func doRef(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("ref (branch, full commit SHA-1 or short commit SHA-1) must be given")
	}
	ref := args[0]

	k8sClient, err := kubernetes.NewClient(rootOpts.annotationPrefix, rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	deployment, err := k8sClient.DetectTargetDeployment(refOpts.namespace, refOpts.deployment)
	if err != nil {
		return errors.Wrap(err, "failed to detect target Deployment")
	}

	container, err := k8sClient.DetectTargetContainer(deployment, refOpts.container)
	if err != nil {
		return errors.Wrap(err, "failed to detect target container")
	}

	repos, err := deployment.Repositories()
	if err != nil {
		return errors.Wrap(err, "failed to extract repositories from deployment")
	}

	repo, ok := repos[container.Name()]
	if !ok {
		return errors.Errorf("GitHub repository for container %q not found in deployment", container.Name())
	}

	ctx := context.Background()
	ghClient := github.NewClient(ctx, refOpts.accessToken)

	sha1, err := ghClient.CommitFronRef(repo, ref)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve commit SHA-1 matched to ref %q in repo %q", ref, repo)
	}

	currentImage := deployment.ContainerImage(container.Name())
	newImage := strings.Split(currentImage, ":")[0] + ":" + sha1

	if refOpts.dryRun {
		fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("[dry-run]  before: %s\n", container.Image())
		fmt.Printf("[dry-run]   after: %s\n", newImage)
	} else {
		fmt.Printf("deploy to (deployment: %q, container: %q)\n", deployment.Name(), container.Name())
		fmt.Printf("  before: %s\n", container.Image())
		fmt.Printf("   after: %s\n", newImage)

		if _, err := k8sClient.SetImage(
			deployment, container.Name(), newImage, refOpts.user, composeRefCause(ref, container.Name(), deployment.Name(), refOpts.namespace),
		); err != nil {
			return errors.Wrap(err, "failed to set image")
		}

		fmt.Printf("\n")
		fmt.Printf("deployment successfully updated! check rollout status by `kubectl rollout status deployment/DEPLOYMENT --namespace %s`\n", refOpts.namespace)
	}

	return nil
}

func composeRefCause(ref, container, deployment, namespace string) string {
	return fmt.Sprintf(`k8ship ref %s --container "%s" --deployment "%s" --namespace "%s"`, ref, container, deployment, namespace)
}

func init() {
	RootCmd.AddCommand(refCmd)

	refCmd.Flags().StringVar(&refOpts.accessToken, "access-token", "", "GitHub access token")
	refCmd.Flags().StringVarP(&refOpts.container, "container", "c", "", "target container")
	refCmd.Flags().StringVarP(&refOpts.deployment, "deployment", "d", "", "target Deployment")
	refCmd.Flags().BoolVar(&refOpts.dryRun, "dry-run", false, "dry run")
	refCmd.Flags().StringVarP(&refOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
	refCmd.Flags().StringVarP(&refOpts.user, "user", "u", "", "image tag to deploy (default: current login user)")

	if refOpts.accessToken == "" {
		refOpts.accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	}

	if refOpts.user == "" {
		refOpts.user = os.Getenv("USER")
	}
}
