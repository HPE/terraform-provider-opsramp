terraform {
  required_providers {
    opsramp = {
      source  = "registry.terraform.io/HPE/opsramp"
      version = ">=0.1.2"
    }
  }
}

provider "opsramp" {
  client_id     = "HDNPG7B8zcDGT2Qg8MnAAWv9srzgRwRz"
  client_secret = "U4gV2UrHhwsme4Hvuf9uEutww5mVg5CBmjaAm9SgR7tJedrmpXTpkc4YKZ9UVQxH"
  endpoint      = "score.api.opsramp.com"
  tenant        = "29f83b5e-f7f0-4cd9-9cdc-53c9677eac87"
}

# Create individual resources
resource "opsramp_resource" "resource1" {
  resource_name = "Test1"
  resource_type = "Server"
}

resource "opsramp_resource" "resource2" {
  resource_name = "Test2"
  resource_type = "Server"
}

# Create a service with multiple resources
resource "opsramp_servicemap" "servicemap_root" {
  name = "TestRoot"
  type = "Service"
}

resource "opsramp_servicemap" "servicemap_child1" {
  name   = "Test1"
  type   = "Service"
  parent = opsramp_servicemap.servicemap_root.id
}

resource "opsramp_servicemap" "servicemap_child2" {
  name   = "Test2"
  type   = "Service"
  parent = opsramp_servicemap.servicemap_root.id
}

resource "opsramp_servicemap" "servicemap_child21" {
  name      = "Test21"
  type      = "Resource"
  parent    = opsramp_servicemap.servicemap_child2.id
  resources = [opsramp_resource.resource1.uuid]
}

resource "opsramp_servicemap" "servicemap_child22" {
  name         = "Test22"
  type         = "Resource"
  parent       = opsramp_servicemap.servicemap_child2.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}

# Service map links
resource "opsramp_servicemap" "servicemap_linked_root" {
  name = "Test Linked Root"
  type = "Service"
}

resource "opsramp_servicemap_link" "servicemap_link" {
  parent = opsramp_servicemap.servicemap_root.id
  link   = opsramp_servicemap.servicemap_linked_root.id
}