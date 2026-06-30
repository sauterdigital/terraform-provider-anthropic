# Federation rule pinning a specific GitHub repo to the ci_deploy service
# account. Subject prefix pattern follows GitHub's OIDC sub format.
resource "anthropic_federation_rule" "github_repo_main" {
  name               = "deploy-from-sauter-website-main"
  issuer_id          = anthropic_federation_issuer.github_actions.id
  service_account_id = anthropic_service_account.ci_deploy.id
  oauth_scope        = "workspace:developer workspace:inference"

  match_subject_prefix = "repo:sauterdigital/website:ref:refs/heads/main"

  workspace_id           = anthropic_workspace.prod.id
  token_lifetime_seconds = 3600
}
