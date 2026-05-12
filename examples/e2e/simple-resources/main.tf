terraform {
  required_providers {
    opsramp = {
      source  = "github.com/HPE/terraform-provider-opsramp"
      version = ">=0.1.0"
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