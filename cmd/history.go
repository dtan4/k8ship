package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

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

		lines := formatHistory(rs)
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		headers := []string{
			"CREATED AT",
			"IMAGES",
		}
		fmt.Fprintln(w, strings.Join(headers, "\t"))

		for _, l := range lines {
			fmt.Fprintln(w, l)
		}

		w.Flush()
	}

	return nil
}

func formatHistory(rs []*kubernetes.ReplicaSet) []string {
	lines := make([]string, 0, len(rs))

	for _, r := range rs {
		lines = append(lines, fmt.Sprintf("%s\t%s", r.CreatedAt(), formatImages(r.Images())))
	}

	return lines
}

func formatImages(images map[string]string) string {
	ss := make([]string, 0, len(images))

	for k, v := range images {
		ss = append(ss, k+" => "+v)
	}

	return strings.Join(ss, ",")
}

func init() {
	RootCmd.AddCommand(historyCmd)

	historyCmd.Flags().StringVarP(&historyOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
