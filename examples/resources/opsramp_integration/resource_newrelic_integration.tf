
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