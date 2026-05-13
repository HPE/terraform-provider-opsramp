terraform {
  required_providers {
    opsramp = {
      source = "registry.terraform.io/HPE/opsramp"
      version = ">=0.1.3"
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
  resource_name = "TestResource1"
  resource_type = "Server"
}

resource "opsramp_resource" "resource2" {
  resource_name = "TestResource2"
  resource_type = "Server"
}