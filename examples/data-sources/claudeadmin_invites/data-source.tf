data "claudeadmin_invites" "all" {}

output "pending_invites" {
  value = [for i in data.claudeadmin_invites.all.invites : i.email if i.status == "pending"]
}
