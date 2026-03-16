package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cristianverduzco/observex/internal/portforward"
)

func newDashboardCmd() *cobra.Command {
	var (
		grafanaPort     int
		prometheusPort  int
		openBrowser     bool
	)

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Open the Grafana dashboard in your browser",
		Long:  `Port-forward Grafana and Prometheus to localhost and open them in your browser.`,
		Example: `  observex dashboard
  observex dashboard --grafana-port 3000 --prometheus-port 9090`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("📊 Starting port-forwards...")
			return portforward.OpenDashboards(kubeconfig, namespace, grafanaPort, prometheusPort, openBrowser)
		},
	}

	cmd.Flags().IntVar(&grafanaPort, "grafana-port", 3000, "Local port for Grafana")
	cmd.Flags().IntVar(&prometheusPort, "prometheus-port", 9090, "Local port for Prometheus")
	cmd.Flags().BoolVar(&openBrowser, "open", true, "Open browser automatically")

	return cmd
}