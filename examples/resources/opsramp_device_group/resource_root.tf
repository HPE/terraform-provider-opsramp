resource "opsramp_device_group" "device_group_root" {
  name      = "Test Group"
  resources = [opsramp_resource.resource1.uuid]
}