package cli

import (
	"github.com/spf13/cobra"
	"github.com/cristianverduzco/observex/internal/status"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check the status of the observability stack",
		Long:  `Check whether Prometheus, Grafana, and Alertmanager are running and healthy.`,
		Example: `  observex status
  observex status --namespace observex-system`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status.Check(kubeconfig, namespace)
		},
	}
}