package cli

import (
	"github.com/spf13/cobra"
)

var (
	kubeconfig string
	namespace  string
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "observex",
		Short: "ObserveX — one-command Kubernetes observability stack",
		Long: `ObserveX bootstraps a production-ready observability stack on any
Kubernetes cluster with a single command.

Deploys Prometheus, Grafana, and Alertmanager pre-configured and
ready to scrape your cluster in under 2 minutes.`,
	}

	root.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig (defaults to in-cluster or $KUBECONFIG)")
	root.PersistentFlags().StringVar(&namespace, "namespace", "observex-system", "Namespace to install the stack into")

	root.AddCommand(newInstallCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newUninstallCmd())
	root.AddCommand(newDashboardCmd())

	return root
}