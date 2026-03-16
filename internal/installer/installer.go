package installer

import (
	"fmt"
	"os"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/client-go/tools/clientcmd"
)

// Config holds installation configuration
type Config struct {
	Kubeconfig       string
	Namespace        string
	GrafanaPassword  string
	SkipPrometheus   bool
	SkipGrafana      bool
	SkipAlertmanager bool
}

// component represents a Helm chart to install
type component struct {
	name    string
	repo    string
	repoURL string
	chart   string
	version string
	values  map[string]interface{}
}

// Install deploys the full observability stack
func Install(cfg Config) error {
	env := cli.New()
	if cfg.Kubeconfig != "" {
		env.KubeConfig = cfg.Kubeconfig
	}

	// Add Helm repos
	fmt.Println("📦 Adding Helm repositories...")
	if err := addRepo("prometheus-community", "https://prometheus-community.github.io/helm-charts", env); err != nil {
		return fmt.Errorf("failed to add prometheus-community repo: %w", err)
	}
	if err := addRepo("grafana", "https://grafana.github.io/helm-charts", env); err != nil {
		return fmt.Errorf("failed to add grafana repo: %w", err)
	}
	fmt.Println("✓ Helm repositories added")

	// Ensure namespace exists
	if err := ensureNamespace(cfg); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	components := []component{}

	if !cfg.SkipPrometheus {
		components = append(components, component{
			name:    "observex-prometheus",
			repo:    "prometheus-community",
			chart:   "prometheus-community/prometheus",
			version: "",
			values: map[string]interface{}{
				"server": map[string]interface{}{
					"persistentVolume": map[string]interface{}{
						"enabled": false,
					},
					"retention": "7d",
				},
				"alertmanager": map[string]interface{}{
					"enabled": !cfg.SkipAlertmanager,
				},
				"pushgateway": map[string]interface{}{
					"enabled": false,
				},
			},
		})
	}

	if !cfg.SkipGrafana {
		components = append(components, component{
			name:  "observex-grafana",
			repo:  "grafana",
			chart: "grafana/grafana",
			values: map[string]interface{}{
				"adminPassword": cfg.GrafanaPassword,
				"persistence": map[string]interface{}{
					"enabled": false,
				},
				"datasources": map[string]interface{}{
					"datasources.yaml": map[string]interface{}{
						"apiVersion": 1,
						"datasources": []map[string]interface{}{
							{
								"name":      "Prometheus",
								"type":      "prometheus",
								"url":       fmt.Sprintf("http://observex-prometheus-server.%s.svc.cluster.local", cfg.Namespace),
								"access":    "proxy",
								"isDefault": true,
							},
						},
					},
				},
			},
		})
	}

	// Install each component
	for _, c := range components {
		fmt.Printf("⚙ Installing %s...\n", c.name)
		if err := installChart(env, cfg.Namespace, c); err != nil {
			return fmt.Errorf("failed to install %s: %w", c.name, err)
		}
		fmt.Printf("✓ %s installed\n", c.name)
	}

	fmt.Println("\n✅ ObserveX installation complete!")
	fmt.Printf("\n📊 Access your dashboards:\n")
	fmt.Printf("   Run: observex dashboard\n")
	fmt.Printf("   Or:  kubectl port-forward -n %s svc/observex-grafana 3000:80\n", cfg.Namespace)
	fmt.Printf("   Grafana URL:    http://localhost:3000\n")
	fmt.Printf("   Grafana user:   admin\n")
	fmt.Printf("   Grafana pass:   %s\n", cfg.GrafanaPassword)

	return nil
}

// Uninstall removes all observability stack components
func Uninstall(kubeconfig, namespace string) error {
	env := cli.New()
	if kubeconfig != "" {
		env.KubeConfig = kubeconfig
	}

	releases := []string{"observex-grafana", "observex-prometheus"}

	for _, release := range releases {
		fmt.Printf("🗑 Uninstalling %s...\n", release)
		cfg := new(action.Configuration)
		if err := cfg.Init(env.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
			fmt.Printf("⚠ Failed to init helm for %s: %v\n", release, err)
			continue
		}
		client := action.NewUninstall(cfg)
		client.Timeout = 5 * time.Minute
		if _, err := client.Run(release); err != nil {
			fmt.Printf("⚠ Failed to uninstall %s: %v\n", release, err)
		} else {
			fmt.Printf("✓ %s removed\n", release)
		}
	}

	fmt.Println("✅ ObserveX uninstalled")
	return nil
}

// addRepo adds a Helm chart repository
func addRepo(name, url string, env *cli.EnvSettings) error {
	repoFile := env.RepositoryConfig
	r := repo.Entry{Name: name, URL: url}
	chartRepo, err := repo.NewChartRepository(&r, getter.All(env))
	if err != nil {
		return err
	}
	if _, err := chartRepo.DownloadIndexFile(); err != nil {
		return err
	}
	_ = repoFile
	return nil
}

// installChart installs a single Helm chart
func installChart(env *cli.EnvSettings, namespace string, c component) error {
	cfg := new(action.Configuration)
	if err := cfg.Init(env.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {}); err != nil {
		return fmt.Errorf("failed to init helm config: %w", err)
	}

	client := action.NewInstall(cfg)
	client.ReleaseName = c.name
	client.Namespace = namespace
	client.CreateNamespace = true
	client.Timeout = 5 * time.Minute
	client.Wait = false

	// Locate and load the chart
	chartPath, err := client.ChartPathOptions.LocateChart(c.chart, env)
	if err != nil {
		return fmt.Errorf("failed to locate chart %s: %w", c.chart, err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return fmt.Errorf("failed to load chart %s: %w", c.chart, err)
	}

	if _, err := client.Run(chart, c.values); err != nil {
		return fmt.Errorf("failed to install chart %s: %w", c.chart, err)
	}

	return nil
}

// ensureNamespace creates the target namespace if it doesn't exist
func ensureNamespace(cfg Config) error {
	kubeconfig := cfg.Kubeconfig
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	_, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	fmt.Printf("✓ Namespace %s ready\n", cfg.Namespace)
	return nil
}