data "claudeadmin_organization" "current" {}

output "org_id" {
  value = data.claudeadmin_organization.current.id
}
