resource "kubernetes_namespace" "monitoring" {
  count = var.enable_monitoring ? 1 : 0

  metadata {
    name = "monitoring"
    labels = {
      "app.kubernetes.io/part-of" = "spaceship-orders"
    }
  }
}

resource "kubernetes_config_map" "grafana_dashboards" {
  count = var.enable_monitoring ? 1 : 0

  metadata {
    name      = "grafana-spaceship-dashboards"
    namespace = kubernetes_namespace.monitoring[0].metadata[0].name
    labels = {
      "app.kubernetes.io/part-of" = "spaceship-orders"
      grafana_dashboard           = "1"
    }
  }

  data = {
    "spaceship-overview.json" = file("${path.module}/../grafana/dashboards/spaceship-overview.json")
    "postgres-exporter.json"  = file("${path.module}/../grafana/dashboards/postgres-exporter.json")
  }
}

resource "helm_release" "vm_single" {
  count = var.enable_monitoring ? 1 : 0

  name             = "vm-single"
  repository       = "https://victoriametrics.github.io/helm-charts/"
  chart            = "victoria-metrics-single"
  namespace        = kubernetes_namespace.monitoring[0].metadata[0].name
  create_namespace = false
  timeout          = 300

  values = [
    file("${path.module}/../k8s/monitoring/vm-single-values.yaml"),
  ]
}

resource "helm_release" "vmalert" {
  count = var.enable_monitoring ? 1 : 0

  name       = "vmalert"
  repository = "https://victoriametrics.github.io/helm-charts/"
  chart      = "victoria-metrics-alert"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  timeout    = 300

  values = [
    file("${path.module}/../k8s/monitoring/vmalert-values.yaml"),
  ]

  depends_on = [helm_release.vm_single]
}

resource "helm_release" "grafana" {
  count = var.enable_monitoring ? 1 : 0

  name       = "grafana"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  namespace  = kubernetes_namespace.monitoring[0].metadata[0].name
  timeout    = 300

  values = [
    file("${path.module}/../k8s/monitoring/grafana-values.yaml"),
  ]

  depends_on = [
    helm_release.vm_single,
    kubernetes_config_map.grafana_dashboards,
  ]
}
