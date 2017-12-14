package cmd

import (
	"fmt"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View deployment history",
	RunE:  doHistory,
}

var historyOpts = struct {
	namespace string
}{}

func doHistory(cmd *cobra.Command, args []string) error {
	client, err := kubernetes.NewClient(rootOpts.annotationPrefix, rootOpts.kubeconfig, rootOpts.context)
	if err != nil {
		return errors.Wrap(err, "failed to create Kubernetes client")
	}

	ds, err := client.ListDeployments(historyOpts.namespace)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Deployments")
	}

	if len(ds) == 0 {
		return errors.Errorf("no Deployment found in namespace %s", historyOpts.namespace)
	}

	for _, d := range ds {
		fmt.Println("===== " + d.Name())

		rs, err := client.ListReplicaSets(d)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve ReplicaSets")
		}

		for _, r := range rs {
			fmt.Printf("%s %s %s\n", r.Revision(), r.Name(), r.CreatedAt())
		}
	}

	return nil
}

func init() {
	RootCmd.AddCommand(historyCmd)

	historyCmd.Flags().StringVarP(&historyOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
