package cmd

import (
	"fmt"
	"os"

	"github.com/dtan4/k8ship/kubernetes"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8ship",
	Short: "Ship image to Kubernetes easily",
}

var rootOpts = struct {
	context    string
	kubeconfig string
}{}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&rootOpts.context, "context", "", "Kubernetes context")
	RootCmd.PersistentFlags().StringVar(&rootOpts.kubeconfig, "kubeconfig", "", "kubeconfig path")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootOpts.kubeconfig == "" {
		if os.Getenv("KUBECONFIG") == "" {
			rootOpts.kubeconfig = kubernetes.DefaultConfigFile()
		} else {
			rootOpts.kubeconfig = os.Getenv("KUBECONFIG")
		}
	}
}
