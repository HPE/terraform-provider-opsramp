resource "opsramp_servicemap" "child_service" {
  name   = "Service Test"
  type   = "Service"
  parent = opsramp_servicemap.root.id
}