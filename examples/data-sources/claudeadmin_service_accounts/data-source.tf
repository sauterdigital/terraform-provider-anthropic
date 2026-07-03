data "claudeadmin_service_accounts" "developers" {
  organization_role = "developer"
  include_archived  = false
}
