resource "claudeadmin_workspace_member" "platform_dev" {
  workspace_id   = claudeadmin_workspace.example.id
  user_id        = data.claudeadmin_organization_member.alice.id
  workspace_role = "workspace_developer"
}
