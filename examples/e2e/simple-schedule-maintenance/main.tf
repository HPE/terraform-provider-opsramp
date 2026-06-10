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

resource "opsramp_device_group" "example_group" {
  name = "Example Device Group"
  resources = [
    opsramp_resource.resource1.uuid,
    opsramp_resource.resource2.uuid
  ]
}

resource "opsramp_schedule_maintenance" "example" {
  name = "Example Maintenance Window"
  description = "This is an example maintenance window created for testing purposes."

  schedule = {
    type = "recurring"
    start_time = "2026-06-10T15:30:00+0100"
    end_time = "2026-06-10T20:00:00+0100"
    end_by = "Never"
    pattern = {
      type = "daily"
      day_frequency = "everyday"
      repeat_frecuency = 2
    }
    timezone = "Europe/London"
    
  }

  run_rba = true
  install_patch = true
  correlate_alerts = true
  run_escalate_action = true
  
  alert_conditions = {
    matching_type = "ALL"
    rules = [
      {
        key = "subject"
        operator = "endswith"
        value = "test"
      },
      {
        key = "description"
        operator = "startswith"
        value = "test"
      },
      {
        key = "serviceName"
        operator = "contains"
        value = "test"
      },
      {
        key = "resourceName"
        operator = "regex"
        value = "test"
      }
    ]
  }

  device_ids = []
  device_group_ids = [opsramp_device_group.example_group.id]
  site_ids = []
  
  notify_before_end_time = "0"
  notify_before_start_time = "0"
  user_ids = []
  user_group_ids = []
}