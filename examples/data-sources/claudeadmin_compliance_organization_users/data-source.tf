data "claudeadmin_compliance_organization_users" "main" {
  organization_id = "org_01ABC..."
}

# Filter SCIM-sourced users
output "scim_users" {
  value = [
    for u in data.claudeadmin_compliance_organization_users.main.users : u
    if u.source_type == "scim"
  ]
}
