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

const (
	defaultHistoryLimit = 10
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View deployment history",
	RunE:  doHistory,
}

var historyOpts = struct {
	all       bool
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

	tds := []*kubernetes.Deployment{}

	for _, d := range ds {
		if d.IsDeployTarget() {
			tds = append(tds, d)
		}
	}

	if len(tds) == 0 {
		return errors.New("no target Deployments found")
	}

	tcs := map[string]*kubernetes.Container{}

	for _, d := range tds {
		c, err := d.DeployTargetContainer()
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve deploy target container of Deployment %q", d.Name())
		}

		tcs[d.Name()] = c
	}

	for _, d := range tds {
		fmt.Println("===== " + d.Name() + " =====")

		rs, err := client.ListReplicaSets(d)
		if err != nil {
			return errors.Wrap(err, "failed to retrieve ReplicaSets")
		}

		lines := formatHistory(rs, tcs[d.Name()])
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))

		if !historyOpts.all && len(lines) > defaultHistoryLimit {
			lines = lines[0:defaultHistoryLimit]
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		headers := []string{
			"DEPLOYED AT",
			"REVISION",
			"USER",
			"IMAGE",
		}
		fmt.Fprintln(w, strings.Join(headers, "\t"))

		for _, l := range lines {
			fmt.Fprintln(w, l)
		}

		w.Flush()

		fmt.Printf("\n")
	}

	return nil
}

func formatHistory(rs []*kubernetes.ReplicaSet, container *kubernetes.Container) []string {
	lines := make([]string, 0, len(rs))

	for _, r := range rs {
		lines = append(lines, strings.Join([]string{r.CreatedAt().String(), r.Revision(), r.DeployUser(), r.Images()[container.Name()]}, "\t"))
	}

	return lines
}

func init() {
	RootCmd.AddCommand(historyCmd)

	historyCmd.Flags().BoolVarP(&historyOpts.all, "all", "a", false, fmt.Sprintf("Print all relases (default: recent %d items)", defaultHistoryLimit))
	historyCmd.Flags().StringVarP(&historyOpts.namespace, "namespace", "n", kubernetes.DefaultNamespace(), "Kubernetes namespace")
}
