resource "claudeadmin_service_account_workspace" "ci_in_prod" {
  service_account_id = claudeadmin_service_account.ci_deploy.id
  workspace_id       = claudeadmin_workspace.prod.id
  workspace_role     = "workspace_developer"
}
