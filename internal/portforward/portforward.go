package portforward

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// OpenDashboards port-forwards Grafana and Prometheus and opens them in the browser
func OpenDashboards(kubeconfig, namespace string, grafanaPort, prometheusPort int, openBrowser bool) error {
	kubeArgs := []string{}
	if kubeconfig != "" {
		kubeArgs = append(kubeArgs, "--kubeconfig", kubeconfig)
	}

	// Port-forward Grafana
	grafanaArgs := append(kubeArgs,
		"port-forward",
		"-n", namespace,
		"svc/observex-grafana",
		fmt.Sprintf("%d:80", grafanaPort),
	)
	grafanaCmd := exec.Command("kubectl", grafanaArgs...)
	grafanaCmd.Stderr = os.Stderr
	if err := grafanaCmd.Start(); err != nil {
		return fmt.Errorf("failed to port-forward grafana: %w", err)
	}
	defer grafanaCmd.Process.Kill()

	// Port-forward Prometheus
	prometheusArgs := append(kubeArgs,
		"port-forward",
		"-n", namespace,
		"svc/observex-prometheus-server",
		fmt.Sprintf("%d:80", prometheusPort),
	)
	prometheusCmd := exec.Command("kubectl", prometheusArgs...)
	prometheusCmd.Stderr = os.Stderr
	if err := prometheusCmd.Start(); err != nil {
		return fmt.Errorf("failed to port-forward prometheus: %w", err)
	}
	defer prometheusCmd.Process.Kill()

	// Wait for Grafana to be ready
	fmt.Printf("⏳ Waiting for Grafana to be ready at http://localhost:%d...\n", grafanaPort)
	grafanaURL := fmt.Sprintf("http://localhost:%d", grafanaPort)
	for i := 0; i < 20; i++ {
		resp, err := http.Get(grafanaURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("\n📊 Dashboards ready:\n")
	fmt.Printf("   Grafana:    http://localhost:%d  (admin / observex-admin)\n", grafanaPort)
	fmt.Printf("   Prometheus: http://localhost:%d\n", prometheusPort)

	if openBrowser {
		openURL(grafanaURL)
	}

	fmt.Println("\nPress Ctrl+C to stop port-forwarding")

	// Block until interrupted
	select {}
}

// openURL opens a URL in the default browser
func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}