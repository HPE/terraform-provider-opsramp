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
  alias_name = "MyFirstResource"
  resource_name = "TestResource1"
  resource_type = "Other"
}

resource "opsramp_resource" "resource2" {
  alias_name = "MySecondResource"
  hostname = "testresource2.local"
  resource_type = "Other"
}