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

# Custom integration with both inbound and outbound
resource "opsramp_integration" "custom_bidirectional" {
  display_name = "Custom Bidirectional"
  application  = "CUSTOM"
  category     = "Custom"

  inbound = {
    auth_type = "OAUTH2"

    map_attributes = [
      {
        entity_type = "INCIDENT"
        third_party_attribute = "description"
        opsramp_attribute = "incident.impact"
        attribute_values = {"JAJA": "TEST","JAJA2": "TEST2"}
      }
    ]

    enable_drop_alerts = true
  }

  outbound = {
    base_uri  = "https://external-system.example.com/webhook"
    auth_type = "OAUTH2"
    username  = "api-user"
    password  = "secret"
    api_key   = "my-api-key"
    api_secret = "my-api-secret"
    grant_type = "PASSWORD"
    access_token_url = "https://external-system.example.com/oauth/token"
    scope     = "alerts"

    additional_properties = {
      custom_key = "custom_value"
    }
  }
}

# Pre-configured integration (e.g. NewRelic) - auth is auto-provisioned on install
resource "opsramp_integration" "newrelic" {
  application  = "NEWRELIC"

  inbound = {
    auth_type = "WEBHOOK"

    map_attributes = [
        {
            third_party_attribute = "device.ip"
	        opsramp_attribute = "alert.alertTime"
	        	        attribute_values = {"JAJA": "TEST","JAJA2": "TEST2"}

	        default_parsing_value = "test"
        	parsing_operators = [
                {
                    operator = "BEFORE"
                    end_word = "IP:"
                }
            ]
        }
    ]

    enable_drop_alerts = false
  }
}


data "opsramp_custom_event_alert_source" "custom_source" {
  name      = "Custom"
}

resource "opsramp_integration" "client_event_integration" {
    display_name = "Custom Integration Example"
    application  = "CUSTOM-EVENT"
    description = "Custom event integration description"
   
    alert_source_id = data.opsramp_custom_event_alert_source.custom_source.id

    inbound = {

        auth_type = "WEBHOOK"

        map_attributes = [
          {
              third_party_attribute = "alert_time"
              opsramp_attribute = "alert.alertTime"
              default_parsing_value = ""
          },
          {
              third_party_attribute = "alert_id"
              opsramp_attribute = "alert.extAlertId"
              parsing_operators = [
                  {
                      operator = "BEFORE"
                      end_word = " is the alert ID"
                  }
              ]
              default_parsing_value = "default_component"
          },
          {
              third_party_attribute = "alert_component"
              opsramp_attribute = "alert.component"
              parsing_operators = [
                  {
                      operator = "AFTER"
                      start_word = "Component:"
                  }
              ]
              default_parsing_value = "default_component"
          },
          {
              third_party_attribute = "alert_resource_ip_address"
              opsramp_attribute = "alert.resource.ipAddress"
              parsing_operators = [
                  {
                      operator = "BETWEEN"
                      start_word = "IPs:"
                      end_word = ","
                  }
              ]
              default_parsing_value = "default_component"
          },
          {
              third_party_attribute = "alert_device_name"
              opsramp_attribute = "alert.deviceName"
              parsing_operators = [
                  {
                      operator = "MATCHES"
                      regex_str = "Resource Name: (.+)"
                  }
              ]
              default_parsing_value = "default_component"
          }
        ]

        enable_drop_alerts = false
        process_definition_ids = ["PROCESS_843e00d7-24d0-4bf1-a705-9d955930f09e"]
    }
}

output "client_event_integration_id" {
    value = opsramp_integration.client_event_integration.id
}
output "client_event_integration_token" {
    value = opsramp_integration.client_event_integration.inbound.token
    sensitive = true
}
