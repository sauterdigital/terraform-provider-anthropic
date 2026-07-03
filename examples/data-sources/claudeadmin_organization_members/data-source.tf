data "claudeadmin_organization_members" "by_email" {
  email = "alice@example.com"
}

data "claudeadmin_organization_members" "all" {}
