package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dtan4/k8ship/github"
	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [BRANCH|COMMIT_SHA1]",
	Short: "Deploy",
	RunE:  doDeploy,
}

var deployOpts = struct {
	accessToken string
	dryRun      bool
	image       string
	namespace   string
	ref         string
	tag         string
	user        string
}{}

func doDeploy(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if deployOpts.image != "" && deployOpts.tag != "" {
			return errors.New("both target image and tag cannot be specified simultaneously")
		}

		if deployOpts.image == "" && deployOpts.tag == "" {
			return errors.New("target image (--image) or tag (--tag) must be specified")
		}
	} else if len(args) == 1 {
		deployOpts.ref = args[0]
	} else {
		return errors.New("--image, --tag, or ref (branch, full commit SHA-1 or short commit SHA-1) must be given")
	}

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

	targetContainers := map[string]*kubernetes.Container{}

	for _, d := range targetDeployments {
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

	var newImage string

	if deployOpts.ref == "" {
		if deployOpts.image != "" {
			newImage = deployOpts.image
		}

		if deployOpts.tag != "" {
			newImage = image + ":" + deployOpts.tag
		}
	} else {
		ctx := context.Background()
		ghClient := github.NewClient(ctx, refOpts.accessToken)

		sha1, err := ghClient.CommitFronRef(repo, deployOpts.ref)
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve commit SHA-1 matched to ref %q in repo %q", deployOpts.ref, repo)
		}

		cc, err := k8sClient.CurrentContext()
		if err != nil {
			return errors.Wrap(err, "failed to retrieve current context")
		}

		did, err := ghClient.CreateDeployment(repo, deployOpts.ref, cc)
		if err != nil {
			return errors.Wrap(err, "failed to create GitHub Deployment")
		}

		fmt.Printf("Deployment ID: %d\n", did)

		newImage = image + ":" + sha1
	}

	if deployOpts.dryRun {
		for _, d := range targetDeployments {
			c := targetContainers[d.Name()]
			fmt.Printf("[dry-run] deploy to (deployment: %q, container: %q)\n", d.Name(), c.Name())
			fmt.Printf("[dry-run]   before: %s\n", c.Image())
			fmt.Printf("[dry-run]   after:  %s\n", newImage)
		}
	} else {
		for _, d := range targetDeployments {
			c := targetContainers[d.Name()]
			fmt.Printf("deploy to (deployment: %q, container: %q)\n", d.Name(), c.Name())
			fmt.Printf("  before: %s\n", c.Image())
			fmt.Printf("  after:  %s\n", newImage)
		}

		for _, d := range targetDeployments {
			c := targetContainers[d.Name()]

			if _, err := k8sClient.SetImage(
				d, c.Name(), newImage, deployOpts.user, composeDeployCause(deployOpts.ref, deployOpts.image, deployOpts.tag, deployOpts.namespace),
			); err != nil {
				return errors.Wrap(err, "failed to set image")
			}
		}

		fmt.Printf("\n")
		fmt.Printf("deployments successfully updated! check rollout status by `kubectl rollout status deployment/DEPLOYMENT --namespace %s`\n", deployOpts.namespace)
	}

	return nil
}

func composeDeployCause(ref, image, tag, namespace string) string {
	if ref != "" {
		return fmt.Sprintf(`k8ship deploy %s --namespace "%s"`, ref, namespace)
	}

	if image != "" {
		return fmt.Sprintf(`k8ship deploy --image %s --namespace "%s"`, image, namespace)
	}

	if tag != "" {
		return fmt.Sprintf(`k8ship deploy --tag %s --namespace "%s"`, tag, namespace)
	}

	return ""
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&deployOpts.accessToken, "access-token", "", "GitHub access token")
	deployCmd.Flags().BoolVar(&deployOpts.dryRun, "dry-run", false, "dry run")
	deployCmd.Flags().StringVar(&deployOpts.image, "image", "", "image to deploy")
	deployCmd.Flags().StringVarP(&deployOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
	deployCmd.Flags().StringVar(&deployOpts.tag, "tag", "", "image tag to deploy")
	deployCmd.Flags().StringVarP(&deployOpts.user, "user", "u", "", "image tag to deploy (default: current login user)")

	if deployOpts.accessToken == "" {
		deployOpts.accessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	}

	if deployOpts.user == "" {
		deployOpts.user = os.Getenv("USER")
	}
}
