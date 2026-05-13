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

resource "opsramp_site" "vmware_site_spain" {
  name    = "Spain Site"
  country = "Spain"
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

resource "opsramp_site" "vmware_site_valencia" {
  parent    = opsramp_site.vmware_site_spain.id
  name      = "Valencia Data Center"
  address   = "Av. del General Avilés, 35-37, Benicalap"
  country   = "Spain"
  zip       = "46035"
  state     = "Comunitat Valenciana"
  city      = "València"
  resources = [opsramp_resource.resource1.uuid]
}

resource "opsramp_site" "vmware_site_madrid" {
  parent       = opsramp_site.vmware_site_spain.id
  name         = "Madrid Data Center"
  address      = "Calle Vicente Aleixandre, 1"
  country      = "Spain"
  zip          = "28232"
  state        = "Madrid"
  city         = "Las Rozas de Madrid"
  search_query = format("uuid = \"%s\"", opsramp_resource.resource2.uuid)
}