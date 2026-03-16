# ObserveX

> A CLI tool that bootstraps a production-ready Prometheus + Grafana + Alertmanager observability stack onto any Kubernetes cluster with a single command.

ObserveX eliminates the manual work of setting up Kubernetes observability. One command deploys and pre-configures the full monitoring stack — Prometheus scraping your cluster, Grafana with Prometheus as a pre-wired datasource, and Alertmanager — ready to use in under 2 minutes.

Built from scratch to demonstrate deep understanding of the Kubernetes observability stack and Go CLI development using the same framework as kubectl and Helm.

---

## Demo
```bash
$ observex install --kubeconfig ~/.kube/config

🚀 ObserveX — Installing observability stack
   Namespace:    observex-system
   Prometheus:   true
   Grafana:      true
   Alertmanager: true

📦 Adding Helm repositories...
✓ Helm repositories added
✓ Namespace observex-system ready
⚙ Installing observex-prometheus...
✓ observex-prometheus installed
⚙ Installing observex-grafana...
✓ observex-grafana installed

✅ ObserveX installation complete!
   Grafana URL:   http://localhost:3000
   Grafana user:  admin
   Grafana pass:  observex-admin
```

---

## Commands

### `observex install`

Deploys Prometheus, Grafana, and Alertmanager to your cluster. Grafana is pre-configured with Prometheus as a datasource — no manual wiring required.
```bash
# Install with defaults
observex install

# Custom Grafana password
observex install --grafana-password mysecretpassword

# Skip Alertmanager
observex install --skip-alertmanager

# Target a specific namespace
observex install --namespace monitoring
```

### `observex status`

Checks whether all components are running and healthy.
```bash
observex status

# Output:
# 🔍 ObserveX Status — namespace: observex-system
#   📈 Prometheus    ✅ Running (observex-prometheus-server-xxx)
#   📊 Grafana       ✅ Running (observex-grafana-xxx)
#   🔔 Alertmanager  ✅ Running (observex-prometheus-alertmanager-0)
# ✅ All components healthy
```

### `observex dashboard`

Port-forwards Grafana and Prometheus to localhost and opens Grafana in your browser.
```bash
observex dashboard

# Custom ports
observex dashboard --grafana-port 3000 --prometheus-port 9090

# Don't auto-open browser
observex dashboard --open=false
```

### `observex uninstall`

Removes all ObserveX components from the cluster.
```bash
# With confirmation prompt
observex uninstall

# Skip confirmation
observex uninstall --force
```

---

## Installation

### Prerequisites

- Go 1.21+
- `kubectl` configured against a Kubernetes cluster
- `helm` installed and repos added:
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

### Build from source
```bash
git clone https://github.com/cristianverduzco/observex
cd observex
go build -o bin/observex ./cmd/observex
```

### Run
```bash
./bin/observex install --kubeconfig ~/.kube/config
```

---

## CLI Flags

### Global flags (all commands)

| Flag | Default | Description |
|---|---|---|
| `--kubeconfig` | `$KUBECONFIG` | Path to kubeconfig file |
| `--namespace` | `observex-system` | Namespace to install into |

### `install` flags

| Flag | Default | Description |
|---|---|---|
| `--grafana-password` | `observex-admin` | Grafana admin password |
| `--skip-prometheus` | `false` | Skip Prometheus installation |
| `--skip-grafana` | `false` | Skip Grafana installation |
| `--skip-alertmanager` | `false` | Skip Alertmanager installation |

### `dashboard` flags

| Flag | Default | Description |
|---|---|---|
| `--grafana-port` | `3000` | Local port for Grafana |
| `--prometheus-port` | `9090` | Local port for Prometheus |
| `--open` | `true` | Auto-open browser |

### `uninstall` flags

| Flag | Default | Description |
|---|---|---|
| `--force` | `false` | Skip confirmation prompt |

---

## What Gets Deployed

| Component | Chart | Description |
|---|---|---|
| Prometheus | `prometheus-community/prometheus` | Metrics collection, scraping, storage |
| Grafana | `grafana/grafana` | Visualization, dashboards, alerting UI |
| Alertmanager | Bundled with Prometheus | Alert routing and notification |
| Node Exporter | Bundled with Prometheus | Host-level metrics (CPU, memory, disk) |
| Kube State Metrics | Bundled with Prometheus | Kubernetes object metrics |

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go |
| CLI framework | Cobra (same as kubectl, Helm) |
| Helm integration | Helm SDK (helm.sh/helm/v3) |
| Kubernetes client | client-go |
| Infrastructure | Kubernetes (kubeadm), Arch Linux |

---

## Roadmap

- [x] `observex install` — full stack deployment
- [x] `observex status` — component health check
- [x] `observex dashboard` — port-forward + browser open
- [x] `observex uninstall` — clean teardown
- [x] Pre-configured Grafana datasource pointing to Prometheus
- [ ] Built-in Grafana dashboards for Kubernetes cluster overview
- [ ] `observex upgrade` — upgrade components to latest versions
- [ ] `observex logs` — tail logs from all components
- [ ] Support for custom Prometheus scrape configs
- [ ] Helm chart for deploying ObserveX itself as a cluster service

---

## Status

✅ Core install, status, dashboard, and uninstall commands complete — tested on a self-hosted kubeadm cluster (Arch Linux, Kubernetes v1.35).