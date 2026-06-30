resource "anthropic_service_account_workspace" "ci_in_prod" {
  service_account_id = anthropic_service_account.ci_deploy.id
  workspace_id       = anthropic_workspace.prod.id
  workspace_role     = "workspace_developer"
}
