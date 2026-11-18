output "app_namespace" {
  value = "spaceship"
}

output "monitoring_enabled" {
  value = var.enable_monitoring
}

output "apply_hint" {
  value = var.enable_monitoring ? "kubectl -n monitoring port-forward svc/grafana 3000:80" : "monitoring skipped (enable_monitoring=false)"
}
