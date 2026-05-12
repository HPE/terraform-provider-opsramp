resource "opsramp_servicemap_link" "servicemap_link" {
  parent = opsramp_servicemap.servicemap_root.id
  link   = opsramp_servicemap.servicemap_linked_root.id
}