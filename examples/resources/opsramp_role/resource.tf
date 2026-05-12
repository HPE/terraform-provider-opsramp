resource "opsramp_role" "client_admin_role" {
  name        = "Admin Role"
  description = "Administrative role with full permissions"

  permissions = [
    opsramp_permission_set.client_admin_perms.unique_id
  ]
}