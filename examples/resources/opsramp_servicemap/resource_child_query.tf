resource "opsramp_servicemap" "child_search" {
  name         = "Search Test"
  type         = "Resource"
  parent       = opsramp_servicemap.root.id
  search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}