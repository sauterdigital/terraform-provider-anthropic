# Service accounts REQUIRE OAuth Bearer auth (oauth_token, not admin_api_key).
resource "anthropic_service_account" "ci_deploy" {
  name              = "ci-deploy-bot"
  description       = "GitHub Actions deploy pipeline for the marketing workspace"
  organization_role = "developer"
}
