# User Group for Client Admins
resource "opsramp_user_group" "client_admin_group" {
  name        = "Admin User Group"
  description = "User group for client administrators"

  roles = [
    opsramp_role.client_admin_role.id
  ]

  users = [
    opsramp_user.admin.id
  ]
}