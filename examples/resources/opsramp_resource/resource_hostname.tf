# Create a resource using hostname and resource_type
resource "opsramp_resource" "test_resource" {
  hostname      = "testresource.example.com"
  resource_type = "Server"
}
