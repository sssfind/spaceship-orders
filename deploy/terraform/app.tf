locals {
  k8s_dir = "${path.module}/../k8s"
  app_manifests = [
    "orders-db.yaml",
    "iam-db.yaml",
    "victoria-logs.yaml",
    "otel-collector.yaml",
    "blackbox.yaml",
    "iam.yaml",
    "order.yaml",
  ]
}

data "kubectl_file_documents" "namespace" {
  content = file("${local.k8s_dir}/namespace.yaml")
}

resource "kubectl_manifest" "namespace" {
  for_each  = data.kubectl_file_documents.namespace.manifests
  yaml_body = each.value
}

data "kubectl_file_documents" "app" {
  for_each = toset(local.app_manifests)
  content  = file("${local.k8s_dir}/${each.value}")
}

locals {
  app_docs = merge([
    for name, data in data.kubectl_file_documents.app : {
      for id, doc in data.manifests : "${name}/${id}" => doc
    }
  ]...)
}

resource "kubectl_manifest" "app" {
  for_each  = local.app_docs
  yaml_body = each.value

  wait              = false
  server_side_apply = false
  validate_schema   = false

  depends_on = [kubectl_manifest.namespace]
}
