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

resource "opsramp_script_category" "automation" {
  name = "Test Automation Scripts"
}

resource "opsramp_script_category" "linux" {
  name      = "Test Linux Scripts"
  parent_id = opsramp_script_category.automation.uuid
}

resource "opsramp_script" "restart_service" {
  category_id     = opsramp_script_category.linux.uuid
  name            = "Restart Service"
  description     = "Restart a service on a Linux machine."
  platforms       = ["LINUX"]
  execution_type  = "SHELL"
  install_timeout = 120

  attachment = {
    name        = "restart_service_linux.sh"
    file = file("./restart_service_linux.sh")
  }

  parameters = [
    {
      name          = "service_name"
      description   = ""
      default_value = ""
      type          = "REQUIRED"
      data_type     = "STRING"
    },
  ]
}