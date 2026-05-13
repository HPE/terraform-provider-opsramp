resource "opsramp_servicemap" "root" {
  name = "Root"
  type = "Service"
}

resource "opsramp_servicemap" "child_resource" {
  name      = "Resource Test"
  type      = "Resource"
  parent    = opsramp_servicemap.root.id
  resources = [opsramp_resource.resource1.uuid]
}

resource "opsramp_servicemap" "child_search" {
  name         = "Search Test"
  type         = "Resource"
  parent       = opsramp_servicemap.root.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}