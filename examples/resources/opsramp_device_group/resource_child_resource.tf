resource "opsramp_device_group" "device_group_resources" {
  parent    = opsramp_device_group.device_group_root.id
  name      = "Test Resources"
  resources = [opsramp_resource.resource2.uuid]
}