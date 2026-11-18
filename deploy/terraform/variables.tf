variable "kubeconfig_path" {
  type        = string
  description = "Kubeconfig (Minikube, kind, или другой кластер)."
  default     = "~/.kube/config"
}

variable "kubeconfig_context" {
  type        = string
  description = "Контекст kubectl. Пустая строка — текущий."
  default     = ""
}

variable "enable_monitoring" {
  type        = bool
  description = "Helm-релизы VictoriaMetrics, vmalert и Grafana в namespace monitoring. В CI обычно false."
  default     = true
}
