resource "opsramp_script_category" "network" {
  name      = "Network"
  parent_id = opsramp_script_category.automation.id
}