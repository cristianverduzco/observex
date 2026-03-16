package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/cristianverduzco/observex/internal/installer"
)

func newUninstallCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the observability stack",
		Long:  `Uninstall Prometheus, Grafana, and Alertmanager from your cluster.`,
		Example: `  observex uninstall
  observex uninstall --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !force {
				fmt.Print("⚠ This will remove all observability components. Continue? [y/N]: ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Aborted.")
					return nil
				}
			}
			return installer.Uninstall(kubeconfig, namespace)
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Skip confirmation prompt")
	return cmd
}