resource "claudeadmin_federation_rule_workspace" "ci_in_staging" {
  federation_rule_id = claudeadmin_federation_rule.github_repo_main.id
  workspace_id       = claudeadmin_workspace.staging.id
}
