package status

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Check prints the status of all observability components
func Check(kubeconfig, namespace string) error {
	client, err := newClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	fmt.Printf("🔍 ObserveX Status — namespace: %s\n\n", namespace)

	components := []struct {
		name    string
		label   string
		emoji   string
	}{
		{"Prometheus", "app.kubernetes.io/name=prometheus", "📈"},
		{"Grafana", "app.kubernetes.io/name=grafana", "📊"},
		{"Alertmanager", "app.kubernetes.io/name=alertmanager", "🔔"},
	}

	allHealthy := true
	for _, c := range components {
		pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
			LabelSelector: c.label,
		})
		if err != nil || len(pods.Items) == 0 {
			fmt.Printf("  %s %-20s ❌ Not found\n", c.emoji, c.name)
			allHealthy = false
			continue
		}

		pod := pods.Items[0]
		ready := false
		for _, cond := range pod.Status.Conditions {
			if string(cond.Type) == "Ready" && string(cond.Status) == "True" {
				ready = true
				break
			}
		}

		if ready {
			fmt.Printf("  %s %-20s ✅ Running (%s)\n", c.emoji, c.name, pod.Name)
		} else {
			fmt.Printf("  %s %-20s ⏳ Starting (%s)\n", c.emoji, c.name, pod.Name)
			allHealthy = false
		}
	}

	fmt.Println()
	if allHealthy {
		fmt.Println("✅ All components healthy")
		fmt.Println("\n💡 Run 'observex dashboard' to open Grafana in your browser")
	} else {
		fmt.Println("⚠ Some components are not ready — run again in a moment")
	}

	return nil
}

func newClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}