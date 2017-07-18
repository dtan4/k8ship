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

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy",
	RunE:  doDeploy,
}

var deployOpts = struct {
	accessToken string
	dryRun      bool
	namespace   string
}{}

func doDeploy(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("ref (branch, full commit SHA-1 or short commit SHA-1) must be given")
	}
	ref := args[0]

	k8sClient, err := kubernetes.NewClient(rootOpts.kubeconfig, rootOpts.context)
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

	if len(deployments) == 0 {
		return errors.New("no target Deployments found")
	}

	targetContainers := map[string]*kubernetes.Container{}

	for _, d := range deployments {
		c, err := d.DeployTargetContainer()
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve deploy target container of Deployment %q", d.Name())
		}

		targetContainers[d.Name()] = c
	}

	repo, err := kubernetes.GetTargetRepository(targetDeployments, targetContainers)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve target repository")
	}

	image, err := kubernetes.GetTargetImage(targetContainers)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve target image")
	}

	ctx := context.Background()
	ghClient := github.NewClient(ctx, refOpts.accessToken)

	sha1, err := ghClient.CommitFronRef(repo, ref)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve commit SHA-1 matched to ref %q in repo %q", ref, repo)
	}

	newImage := strings.Split(image, ":")[0] + ":" + sha1

	if deployOpts.dryRun {
		for _, d := range targetDeployments {
			fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", d.Name(), targetContainers[d.Name()].Name())
		}
		fmt.Printf("[dry-run]  before: %s\n", image)
		fmt.Printf("[dry-run]   after: %s\n", newImage)
	} else {
		for _, d := range targetDeployments {
			fmt.Printf("deploy to (deployment: %q, container: %q)\n", d.Name(), targetContainers[d.Name()].Name())
		}
		fmt.Printf("  before: %s\n", image)
		fmt.Printf("   after: %s\n", newImage)

		for _, d := range targetDeployments {
			c := targetContainers[d.Name()]

			if _, err := k8sClient.SetImage(d, c.Name(), newImage); err != nil {
				return errors.Wrap(err, "failed to set image")
			}
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&deployOpts.accessToken, "access-token", "", "GitHub access token")
	deployCmd.Flags().BoolVar(&deployOpts.dryRun, "dry-run", false, "dry run")
	deployCmd.Flags().StringVarP(&deployOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")

	if deployOpts.accessToken == "" {
		deployOpts.accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	}
}
