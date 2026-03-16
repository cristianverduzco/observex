package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cristianverduzco/observex/internal/installer"
)

func newInstallCmd() *cobra.Command {
	var (
		grafanaPassword string
		skipPrometheus  bool
		skipGrafana     bool
		skipAlertmanager bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the full observability stack",
		Long: `Install Prometheus, Grafana, and Alertmanager on your Kubernetes cluster.

All components are pre-configured and ready to scrape your cluster
immediately after installation.`,
		Example: `  # Install with default settings
  observex install

  # Install with custom Grafana password
  observex install --grafana-password mysecretpassword

  # Install only Prometheus and Grafana
  observex install --skip-alertmanager`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := installer.Config{
				Kubeconfig:       kubeconfig,
				Namespace:        namespace,
				GrafanaPassword:  grafanaPassword,
				SkipPrometheus:   skipPrometheus,
				SkipGrafana:      skipGrafana,
				SkipAlertmanager: skipAlertmanager,
			}

			fmt.Println("🚀 ObserveX — Installing observability stack")
			fmt.Printf("   Namespace:  %s\n", namespace)
			fmt.Printf("   Prometheus: %v\n", !skipPrometheus)
			fmt.Printf("   Grafana:    %v\n", !skipGrafana)
			fmt.Printf("   Alertmanager: %v\n\n", !skipAlertmanager)

			return installer.Install(cfg)
		},
	}

	cmd.Flags().StringVar(&grafanaPassword, "grafana-password", "observex-admin", "Grafana admin password")
	cmd.Flags().BoolVar(&skipPrometheus, "skip-prometheus", false, "Skip Prometheus installation")
	cmd.Flags().BoolVar(&skipGrafana, "skip-grafana", false, "Skip Grafana installation")
	cmd.Flags().BoolVar(&skipAlertmanager, "skip-alertmanager", false, "Skip Alertmanager installation")

	return cmd
}