data "claudeadmin_compliance_group_members" "engineering" {
  group_id = "group_01XYZ..."
}

output "engineering_emails" {
  value = [for m in data.claudeadmin_compliance_group_members.engineering.members : m.email]
}
