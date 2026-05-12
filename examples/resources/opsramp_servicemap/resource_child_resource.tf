resource "opsramp_servicemap" "child_resource" {
  name      = "Resource Test"
  type      = "Resource"
  parent    = opsramp_servicemap.root.id
  resources = [opsramp_resource.resource1.uuid]
}