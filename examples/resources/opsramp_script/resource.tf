resource "opsramp_script" "restart_service" {
  category_id     = opsramp_script_category.os_management.id
  name            = "Restart Service"
  description     = "Restarts a named OS service."
  platforms       = ["LINUX"]
  execution_type  = "SHELL"
  install_timeout = 120
  service_name    = "opsramp-agent"

  attachment = {
    name        = "restart_service_linux.sh"
    content_url = base64encode(file("restart_service_linux.sh"))
  }

  parameters = [
    {
      name          = "service_name"
      description   = "Name of the service to restart"
      default_value = "opsramp-agent"
      type          = "REQUIRED"
      data_type     = "STRING"
    },
  ]
}