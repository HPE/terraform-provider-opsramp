terraform {
  required_providers {
    opsramp = {
      source = "registry.terraform.io/HPE/opsramp"
      version = ">=0.1.4"
    }
  }
}

provider "opsramp" {
  client_id     = "*****"
  client_secret = "*****"
  endpoint      = "*****"
  tenant        = "*****"
}

# Create individual resources
resource "opsramp_resource" "resource1" {
  resource_name = "Test1"
  resource_type = "Linux"
}

resource "opsramp_resource" "resource2" {
  resource_name = "Test2"
  resource_type = "Linux"
}

resource "opsramp_resource" "resource3" {
  resource_name = "Test3"
  resource_type = "Linux"
}

resource "opsramp_device_group" "device_group_root" {
  name      = "Test Groups"
  resources = []
}

resource "opsramp_device_group" "device_group_resources" {
  parent_id = opsramp_device_group.device_group_root.id
  name      = "Test Resources"
  resources = [opsramp_resource.resource1.uuid]
}

resource "opsramp_device_group" "device_group_query" {
  parent_id    = opsramp_device_group.device_group_root.id
  name         = "Test Queries"
  search_query = format("resourceType = \"Linux\" AND uuid = \"%s\"", opsramp_resource.resource2.uuid)
}

resource "opsramp_device_group" "device_group_mixed" {
  parent_id    = opsramp_device_group.device_group_root.id
  name         = "Test Queries Mixed"
  search_query = format("resourceType = \"Linux\" AND uuid = \"%s\"", opsramp_resource.resource2.uuid)
  resources    = [opsramp_resource.resource3.uuid]
}
