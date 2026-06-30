resource "anthropic_federation_rule_workspace" "ci_in_staging" {
  federation_rule_id = anthropic_federation_rule.github_repo_main.id
  workspace_id       = anthropic_workspace.staging.id
}
