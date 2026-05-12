# Create a resource using resource_name and resource_type
resource "opsramp_resource" "test_resource" {
  resource_name = "TestResource"
  resource_type = "Server"
}